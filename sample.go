// Copyright 2015 someonegg. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package samplefactory

import (
	"fmt"

	"os"
	"sort"

	"encoding/csv"

	"path/filepath"

	"github.com/zxfonline/csvconfig"
	"github.com/zxfonline/fileutil"
	"github.com/zxfonline/go-ulua"
	"github.com/zxfonline/golog"
	"github.com/zxfonline/json"
	"github.com/zxfonline/xlsx"
)

var (
	logger = golog.New("Sample")
)

//模板接口
type Sample interface {
	//模板sid
	GetSid() int
	//数据克隆
	Clone() Sample
}
type SampleFactory struct {
	samples   map[int]Sample
	tableName string
}

//获取模板数据
func (f *SampleFactory) GetSample(sid int) Sample {
	return f.samples[sid].Clone()
}

//获取模板克隆数据
func (f *SampleFactory) NewSample(sid int) Sample {
	return f.samples[sid].Clone()
}

//解析csv文件生成模板对象工厂
func CreateSampleFacotry(tableName string, nf func() Sample) (*SampleFactory, error) {
	rc := csvconfig.GetAll(tableName)
	if rc == nil {
		return nil, fmt.Errorf("config no found,check the csv[%v] config is initialized", tableName)
	}
	factory := &SampleFactory{samples: make(map[int]Sample), tableName: tableName}
	for _, tb := range rc {
		fields, err := json.Marshal(tb.Fields)
		if err != nil {
			return nil, err
		}
		sample := nf()
		err = json.Unmarshal(fields, &sample)
		if err != nil {
			return nil, err
		}
		factory.samples[sample.GetSid()] = sample
	}
	return factory, nil
}

//数据本地化 saveType：1=excel、2=csv、4=lua
func (f *SampleFactory) Store(savePath string, saveType int, nf func() Sample) error {
	savePath = filepath.ToSlash(savePath)
	if saveType&1 != 0 || saveType&2 != 0 {
		objptr := nf()
		toj, err := json.Marshal(&objptr)
		if err != nil {
			return err
		}
		sampleTitle := make(map[string]string)
		err = json.Unmarshal(toj, &sampleTitle)
		if err != nil {
			return err
		}
		record := len(f.samples) + 1
		fields := len(sampleTitle)
		table := make([][]string, record, record)

		//title构建
		titles := make([]string, fields, fields)
		table[0] = titles
		idx := 0
		for k, _ := range sampleTitle {
			titles[idx] = k
			idx++
		}

		//sid排序
		recodes := make([]Sample, 0, len(f.samples))
		for _, sample := range f.samples {
			recodes = append(recodes, sample)
		}
		sort.Sort(recordHeap(recodes))

		//content构建
		idx = 1
		for _, sample := range recodes {
			toj, err = json.Marshal(&sample)
			if err != nil {
				return err
			}
			sampleTitle = make(map[string]string)
			err = json.Unmarshal(toj, &sampleTitle)
			if err != nil {
				return err
			}
			table[idx] = make([]string, fields, fields)
			for ti, tk := range titles {
				table[idx][ti] = sampleTitle[tk]
			}
			idx++
		}
		if saveType&1 != 0 {
			err = f.saveExcel(table, savePath)
			if err != nil {
				return err
			}
		}
		if saveType&2 != 0 {
			err = f.saveCsv(table, savePath)
			if err != nil {
				return err
			}
		}
	} else if saveType&4 != 0 {
		err := f.saveLua(savePath)
		if err != nil {
			return err
		}
	}
	return nil
}
func (f *SampleFactory) saveExcel(table [][]string, savePath string) (err error) {
	xlsxFile := xlsx.NewFile()
	var sheet *xlsx.Sheet
	sheet, err = xlsxFile.AddSheet(f.tableName)
	if err != nil {
		return
	}
	for _, record := range table {
		row := sheet.AddRow()
		for _, field := range record {
			cell := row.AddCell()
			cell.Value = field
		}
	}
	savePath = fileutil.PathJoin(savePath, f.tableName+".xlsx")
	if fileutil.FileExists(savePath) {
		logger.Warnf("store excel file exists.rewrite path=%v", savePath)
	}
	var file *os.File
	file, err = fileutil.OpenFile(savePath, fileutil.DefaultFileMode, os.O_TRUNC|os.O_CREATE|os.O_WRONLY)
	if err != nil {
		return
	}
	err = xlsxFile.Write(file)
	if err != nil {
		return
	}
	return file.Close()
}
func (f *SampleFactory) saveCsv(table [][]string, savePath string) (err error) {
	savePath = fileutil.PathJoin(savePath, f.tableName+".csv")
	if fileutil.FileExists(savePath) {
		logger.Warnf("store csv file exists.rewrite path=%v", savePath)
	}
	var file *os.File
	file, err = fileutil.OpenFile(savePath, fileutil.DefaultFileMode, os.O_TRUNC|os.O_CREATE|os.O_WRONLY)
	if err != nil {
		return
	}
	defer file.Close()
	w := csv.NewWriter(file)
	for _, record := range table {
		if err := w.Write(record); err != nil {
			return fmt.Errorf("error writing record to csv:%v", err)
		}
	}
	w.Flush()
	return w.Error()
}

func (f *SampleFactory) saveLua(savePath string) (err error) {
	savePath = fileutil.PathJoin(savePath, f.tableName+".lua")
	if fileutil.FileExists(savePath) {
		logger.Warnf("store lua file exists.rewrite path=%v", savePath)
	}
	var bb []byte
	bb, err = ulua.MarshalIndent(f.samples, "\t")
	if err != nil {
		return
	}
	var file *os.File
	file, err = fileutil.OpenFile(savePath, fileutil.DefaultFileMode, os.O_TRUNC|os.O_CREATE|os.O_WRONLY)
	if err != nil {
		return
	}
	defer file.Close()
	file.WriteString(f.tableName)
	file.WriteString("=")
	file.Write(bb)
	return
}

type recordHeap []Sample

func (h recordHeap) Len() int {
	return len(h)
}
func (h recordHeap) Less(i, j int) bool {
	return h[i].GetSid() < h[j].GetSid()
}
func (h recordHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}
