package main

import (
	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/analysis/analyzer/keyword"
	"github.com/blevesearch/bleve/analysis/lang/en"
	"github.com/blevesearch/bleve/mapping"
	"log"
)

var (
	indexPath = "logsvc.bleve"
)

func buildIndexMapping() (mapping.IndexMapping, error) {

	// a generic reusable mapping for english text
	englishTextFieldMapping := bleve.NewTextFieldMapping()
	englishTextFieldMapping.Analyzer = en.AnalyzerName

	// a generic reusable mapping for keyword text
	keywordFieldMapping := bleve.NewTextFieldMapping()
	keywordFieldMapping.Analyzer = keyword.Name

	logMapping := bleve.NewDocumentMapping()

	logMapping.AddFieldMappingsAt("Msg", englishTextFieldMapping)
	logMapping.AddFieldMappingsAt("App", keywordFieldMapping)
	logMapping.AddFieldMappingsAt("Tag", keywordFieldMapping)

	indexMapping := bleve.NewIndexMapping()
	indexMapping.AddDocumentMapping("log", logMapping)

	indexMapping.TypeField = "type"
	indexMapping.DefaultAnalyzer = "en"

	return indexMapping, nil
}

func searchApi() {

}

func searcherInit() bleve.Index {
	logIndex, err := bleve.Open(indexPath)
	if err == bleve.ErrorIndexPathDoesNotExist {
		indexMapping, err := buildIndexMapping()
		if err != nil {
			log.Fatalln("srchengInit buildIndexMapping error", err)
		}
		logIndex, err = bleve.New(indexPath, indexMapping)
		if err != nil {
			log.Fatal("srchengInit bleve.New error", err)
		}
	} else if err != nil {
		log.Fatalln("srchengInit bleve.Open error", err)
	}
	return logIndex
}
