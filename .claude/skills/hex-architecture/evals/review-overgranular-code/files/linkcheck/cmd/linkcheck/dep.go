package main

import (
	filelinkstrg "example.com/linkcheck/internal/driven/storage/file/linkstrg"
	filereportstrg "example.com/linkcheck/internal/driven/storage/file/reportstrg"
	filestatestrg "example.com/linkcheck/internal/driven/storage/file/statestrg"
)

func initDeps(cfg config) (*deps, error) {
	links, err := filelinkstrg.New(filelinkstrg.Config{Path: cfg.SourcesFile})
	if err != nil {
		return nil, err
	}
	report, err := filereportstrg.New(filereportstrg.Config{Path: cfg.ReportFile})
	if err != nil {
		return nil, err
	}
	state, err := filestatestrg.New(filestatestrg.Config{Path: cfg.StateFile})
	if err != nil {
		return nil, err
	}
	return &deps{links: links, report: report, state: state}, nil
}
