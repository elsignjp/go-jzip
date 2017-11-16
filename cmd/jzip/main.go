package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/elsignjp/go-jzip"

	"github.com/urfave/cli"
	"gopkg.in/cheggaaa/pb.v1"
)

func main() {
	app := cli.NewApp()
	app.Name = "jzip"
	app.HelpName = "jzip"
	app.Usage = "will convert zip codes of Japan Post"
	app.UsageText = "jzip [global options]"
	app.Version = "0.0.1"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "config, c",
			Value: "",
			Usage: "Specify the config file to be read",
		},
		cli.BoolFlag{
			Name:  "json, j",
			Usage: "Output in JSON format",
		},
		cli.BoolFlag{
			Name:  "sql, q",
			Usage: "Output in SQL format",
		},
		cli.StringFlag{
			Name:  "file, f",
			Value: "",
			Usage: "Specify the ZIP file to be read",
		},
		cli.StringFlag{
			Name:  "out, o",
			Value: "",
			Usage: "Specify the file path to output",
		},
		cli.StringFlag{
			Name:  "table",
			Value: "",
			Usage: "Specify the table name to output (format:SQL)",
		},
		cli.StringFlag{
			Name:  "fld-local",
			Value: "",
			Usage: "Field name of the local public organization code to be output",
		},
		cli.StringFlag{
			Name:  "fld-zip",
			Value: "",
			Usage: "Field name of the zip code to be output",
		},
		cli.StringFlag{
			Name:  "fld-pref",
			Value: "",
			Usage: "Field name of prefecture to output",
		},
		cli.StringFlag{
			Name:  "fld-city",
			Value: "",
			Usage: "Field name of city to output",
		},
		cli.StringFlag{
			Name:  "fld-town",
			Value: "",
			Usage: "Field name of the town area to output",
		},
	}
	app.Action = func(c *cli.Context) error {
		cpath := c.GlobalString("config")
		cfg := initConfig(cpath)
		{
			if c.GlobalBool("json") {
				cfg.Format = "json"
			} else if c.GlobalBool("sql") {
				cfg.Format = "sql"
			}
			if of := c.GlobalString("file"); of != "" {
				cfg.ZipFile = of
			}
			if oo := c.GlobalString("out"); oo != "" {
				cfg.Output = oo
			}
			if ota := c.GlobalString("table"); ota != "" {
				cfg.TableName = ota
			}
			if olo := c.GlobalString("fld-local"); olo != "" {
				cfg.LocalCodeField = olo
			}
			if ozi := c.GlobalString("fld-zip"); ozi != "" {
				cfg.ZipCodeField = ozi
			}
			if opr := c.GlobalString("fld-pref"); opr != "" {
				cfg.PrefecturesField = opr
			}
			if oci := c.GlobalString("fld-city"); oci != "" {
				cfg.CityField = oci
			}
			if oto := c.GlobalString("fld-town"); oto != "" {
				cfg.TownField = oto
			}
		}
		bar := pb.StartNew(0)
		bar.RefreshRate = 200 * time.Millisecond
		bar.AlwaysUpdate = true
		defer bar.Finish()
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		var wr jzip.Writer
		fp, err := os.OpenFile(absolutePath(cfg.Output), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0777)
		defer fp.Close()
		if err != nil {
			fmt.Println("error:", err)
			return nil
		}
		wrcfg := jzip.WriterConfig{
			TableName:        cfg.TableName,
			LocalCodeField:   cfg.LocalCodeField,
			ZipCodeField:     cfg.ZipCodeField,
			PrefecturesField: cfg.PrefecturesField,
			CityField:        cfg.CityField,
			TownField:        cfg.TownField,
		}
		switch cfg.Format {
		case "sql":
			wr = jzip.NewWriter(jzip.SQLWriter, fp, wrcfg)
		default:
			wr = jzip.NewWriter(jzip.JSONWriter, fp, wrcfg)
		}
		rCH, err := jzip.Extract(ctx, cfg.ZipFile)
		if err != nil {
			fmt.Println("error:", err)
			return nil
		}
		sCH := make(chan os.Signal, 1)
		signal.Notify(sCH, syscall.SIGTERM, syscall.SIGINT, os.Interrupt)
		defer signal.Stop(sCH)
		go func() {
			<-sCH
			cancel()
			fmt.Println("-- canceled.")
		}()
		var mu sync.Mutex
		wg := sync.WaitGroup{}
		wg.Add(2)
		wrCH := make(chan []string, 1000)
		go func() {
			defer wg.Done()
			defer close(wrCH)
			for {
				select {
				case <-ctx.Done():
					return
				case r, ok := <-rCH:
					if !ok {
						return
					}
					wrCH <- r
					bar.Total++
				}
			}
		}()
		go func() {
			defer wg.Done()
			nm := jzip.NewNormalize()
			for {
				select {
				case <-ctx.Done():
					return
				case r, ok := <-wrCH:
					if !ok {
						return
					}
					wg.Add(1)
					mu.Lock()
					rs := nm.Add(r)
					if rs != nil {
						for _, nr := range rs {
							if err := wr.Write(nr); err != nil {
								fmt.Println("error:", err)
							}
						}
					}
					bar.Add(1)
					mu.Unlock()
					wg.Done()
				}
			}
		}()
		wg.Wait()
		wr.Close()
		return nil
	}
	app.Run(os.Args)
}
