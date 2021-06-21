package cmd

import (
	"catalogue/api/internal/app"
	"catalogue/api/internal/config"
	"catalogue/api/internal/core/data"
	"catalogue/api/internal/core/service"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/spf13/cobra"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
	"gopkg.in/DataDog/dd-trace-go.v1/profiler"
)

var migDir = "./internal/core/data/migrations"

var rootCmd = &cobra.Command{
	Use:   "serve",
	Short: "Catalogue is used to maintain the product catalogue of an ecommerce company.",
	Long:  `Catalogue is used to maintain the product catalogue of an ecommerce company. It has a set of Micro Services (REST APIs with JSON format responses). `,
	Run: func(cmd *cobra.Command, args []string) {
		runAPIServer()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func runAPIServer() {

	handler := setUpServer()

	if handler.DatadogEnableAPM {
		// datadog tracer
		tracer.Start(
			tracer.WithAgentAddr("dd-agent"),
			tracer.WithServiceName(handler.DatadogServiceName),
			tracer.WithEnv("debug"),
			tracer.WithAnalytics(true),
		)
		defer tracer.Stop()

		if handler.DatadogEnableProfiler {
			// datadog profiler
			err := profiler.Start(
				profiler.WithAgentAddr("dd-agent"),
				profiler.WithService(handler.DatadogServiceName),
				profiler.WithVersion("1.0"),
				profiler.WithTags(fmt.Sprintf("version:%s", "1.0")),
				profiler.WithEnv("debug"),
				profiler.WithProfileTypes(
					profiler.CPUProfile,
					profiler.HeapProfile,
					profiler.GoroutineProfile,
				),
			)
			if err != nil {
				log.Fatal(fmt.Sprintf("failed to start datadog profiler: %v", err))
			}
			defer profiler.Stop()
		}

	}
	router := handler.InitRouter()
	log.Fatal(http.ListenAndServe(handler.Port, router))
}

func setUpServer() *app.Handler {
	handler := new(app.Handler)
	conf := config.NewConfig()
	handler.Config = conf
	dbConn := handler.GetDBConnectionString()
	db, err := handler.GetDBConnection(dbConn)
	if err != nil {
		log.Fatal(err)
	}
	handler.MigrateFolder(dbConn, migDir)
	categoryRepo := data.NewCategoryRepo(db)
	productRepo := data.NewProductRepo(db)
	variantRepo := data.NewVariantRepo(db)
	handler.Category = service.NewCategoryService(categoryRepo)
	handler.Product = service.NewProductService(productRepo)
	handler.Variant = service.NewVariantService(variantRepo)
	return handler
}
