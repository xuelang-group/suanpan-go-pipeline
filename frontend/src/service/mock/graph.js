export default { 
  "htime": 1637566353509, 
  "scale": 1, 
  "x": 0, 
  "y": 0, 
  "nodes": [
    { 
      "uuid": "52fb96d04b6611ec8460c5b86d0280d2", 
      "type": "dataset", 
      "key": "CsvUploader", 
      "name": "csv上传", 
      "x": 493, "y": 162 
    }, 
    { 
      "uuid": "536332404b6611ec8460c5b86d0280d2", 
      "type": "preprocess", 
      "key": "DataSpliter", 
      "name": "拆分训练集", 
      "x": 243, 
      "y": 352 
    }, 
    { 
      "uuid": "53fd28504b6611ec8460c5b86d0280d2", 
      "type": "featureEngineering", 
      "key": "NormalizeComponent", 
      "name": "归一化", 
      "x": 635, 
      "y": 333 
    }, 
    { 
      "uuid": "547b0ae04b6611ec8460c5b86d0280d2", 
      "type": "featureEngineering", 
      "key": "ScaleComponent", 
      "name": "标准归一化", 
      "x": 513, 
      "y": 528 
    }], 
    "connectors": [
      { 
        "tgt": { 
          "uuid": "536332404b6611ec8460c5b86d0280d2", 
          "name": "拆分训练集", 
          "port": "in1" 
        }, 
        "src": { 
          "uuid": "52fb96d04b6611ec8460c5b86d0280d2", 
          "name": "csv上传", 
          "port": "out1" 
        } 
      }, 
      { 
        "tgt": { 
          "uuid": "53fd28504b6611ec8460c5b86d0280d2", 
          "name": "归一化", 
          "port": "in1" 
        }, 
        "src": { 
          "uuid": "52fb96d04b6611ec8460c5b86d0280d2", 
          "name": "csv上传", 
          "port": "out1" 
        } 
      }, 
      { 
        "tgt": { 
          "uuid": "547b0ae04b6611ec8460c5b86d0280d2", 
          "name": "标准归一化", 
          "port": "in1" 
        }, 
        "src": { 
          "uuid": "53fd28504b6611ec8460c5b86d0280d2", 
          "name": "归一化", 
          "port": "out1" 
        } 
      }
    ] 
  }