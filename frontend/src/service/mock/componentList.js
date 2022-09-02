export default [
  {
    "type": "dataProcess",
    "typeLabel": "数据处理",
    "category": null,
    "categoryLabel": null,
    "name": "数据处理",
    "key": "DataProcess",
    "supportConfigInPorts": true,
    "supportConfigOutPorts": true,
    "helpUrl": null,
    "parameters": [
      {
        "key": "expressions",
        "name": "表达式",
        "type": "inputString",
        "required": true
      },
      {
        "key": "inPorts",
        "name": "输入端口",
        "type": "ports",
        "default": [
          {
            "id": "in1",
            "name": "in1"
          }
        ]
      },
      {
        "key": "outPorts",
        "name": "输出端口",
        "type": "ports",
        "default": [
          {
            "id": "out1",
            "name": "out1"
          }
        ]
      }
    ]
  },
  {
    "type": "database",
    "typeLabel": "数据库",
    "category": "read",
    "categoryLabel": "数据库读取",
    "name": "Postgres数据库读取",
    "key": "PostgresReader",
    "helpUrl": null,
    "parameters": [
      {
        "key": "host",
        "name": "主机名或IP地址",
        "type": "inputString",
        "required": true
      },
      {
        "key": "port",
        "name": "端口",
        "type": "inputInteger",
        "required": true
      },
      {
        "key": "user",
        "name": "用户名",
        "type": "inputString",
        "required": true
      },
      {
        "key": "password",
        "name": "密码",
        "type": "inputString",
        "required": true
      },
      {
        "key": "database",
        "name": "数据库",
        "type": "inputString",
        "default": "postgres"
      },
      {
        "key": "sql",
        "name": "查询语句",
        "type": "inputString",
        "default": ""
      }
    ],
    "ports": {
      "in": [
        {
          "id": "in1",
          "name": "in1"
        }
      ],
      "out": [
        {
          "id": "out1",
          "name": "out1"
        }
      ]
    }
  },
  {
    "type": "database",
    "typeLabel": "数据库",
    "category": "write",
    "categoryLabel": "数据库写入",
    "name": "Postgres数据库写入",
    "key": "PostgresWriter",
    "helpUrl": null,
    "parameters": [
      {
        "key": "host",
        "name": "主机名或IP地址",
        "type": "inputString",
        "required": true
      },
      {
        "key": "port",
        "name": "端口",
        "type": "inputInteger",
        "required": true
      },
      {
        "key": "user",
        "name": "用户名",
        "type": "inputString",
        "required": true
      },
      {
        "key": "password",
        "name": "密码",
        "type": "inputString",
        "required": true
      },
      {
        "key": "database",
        "name": "连接数据库",
        "type": "inputString",
        "default": "postgres"
      },
      {
        "key": "databaseChoose",
        "name": "写入数据库",
        "type": "inputString",
        "default": "postgres"
      },
      {
        "key": "table",
        "name": "写入表名",
        "type": "inputString",
        "required": true
      },
      {
        "key": "mode",
        "name": "写入方式",
        "options": [
          {
            "value": "append",
            "label": "追加"
          },
          {
            "value": "replace",
            "label": "覆盖"
          }
        ],
        "type": "select",
        "default": "append"
      },
      {
        "key": "chunksize",
        "name": "每次写入最大行数",
        "type": "inputInteger",
        "default": 100000
      }
    ],
    "ports": {
      "in": [
        {
          "id": "in1",
          "name": "in1"
        }
      ],
      "out": [
        {
          "id": "out1",
          "name": "out1"
        }
      ]
    }
  },
  {
    "type": "dataSimulator",
    "typeLabel": "基础类型",
    "category": null,
    "categoryLabel": null,
    "name": "字符串生成",
    "key": "StringGenerator",
    "helpUrl": null,
    "parameters": [
      {
        "key": "value",
        "name": "字符串值",
        "type": "inputString",
        "required": true
      }
    ],
    "ports": {
      "in": [
        {
          "id": "in1",
          "name": "in1"
        }
      ],
      "out": [
        {
          "id": "out1",
          "name": "out1"
        }
      ]
    }
  },
  {
    "type": "loop",
    "typeLabel": "逻辑语句",
    "category": null,
    "categoryLabel": null,
    "name": "For（触发循环）",
    "key": "ForLoop",
    "helpUrl": null,
    "parameters": [
      {
        "key": "iterator",
        "name": "迭代器",
        "type": "inputString",
        "required": true
      }
    ],
    "ports": {
      "in": [
        {
          "id": "in1",
          "name": "in1"
        }
      ],
      "out": [
        {
          "id": "out1",
          "name": "out1"
        }
      ]
    }
  },
  {
    "type": "script",
    "typeLabel": "代码编辑器",
    "category": null,
    "categoryLabel": null,
    "name": "Python脚本编辑器",
    "key": "ExecutePythonScript",
    "supportConfigInPorts": true,
    "supportConfigOutPorts": true,
    "helpUrl": null,
    "parameters": [
      {
        "key": "script",
        "name": "Python脚本",
        "type": "inputPythonScript",
        "required": true,
        "default": "# The script MUST contain a function named run\\n# which is the entry point for this module.\\n# The entry point function can contain several input arguments:\\n#   Param<in1>: a pandas.DataFrame\\n#   Param<in2>: a pandas.DataFrame\\ndef run(in1 = None, in2 = None):\\n    return in1+in2, in1-in2\\n"
      },
      {
        "key": "inPorts",
        "name": "输入端口",
        "type": "ports",
        "default": [
          {
            "id": "in1",
            "name": "in1"
          },
          {
            "id": "in2",
            "name": "in2"
          }
        ]
      },
      {
        "key": "outPorts",
        "name": "输出端口",
        "type": "ports",
        "default": [
          {
            "id": "out1",
            "name": "out1"
          },
          {
            "id": "out2",
            "name": "out2"
          }
        ]
      }
    ]
  },
  {
    "type": "statement",
    "typeLabel": "逻辑语句",
    "category": null,
    "categoryLabel": null,
    "name": "If语句",
    "key": "IfCondition",
    "helpUrl": null,
    "parameters": [
      {
        "key": "condition",
        "name": "条件判断",
        "type": "inputString",
        "required": true
      }
    ],
    "ports": {
      "in": [
        {
          "id": "in1",
          "name": "in1"
        }
      ],
      "out": [
        {
          "id": "out1",
          "name": "out1"
        },
        {
          "id": "out2",
          "name": "out2"
        }
      ]
    }
  },
  {
    "type": "statement",
    "typeLabel": "逻辑语句",
    "category": null,
    "categoryLabel": null,
    "name": "设置全局变量",
    "key": "GlobalVarsSetting",
    "helpUrl": null,
    "parameters": [
      {
        "key": "name",
        "name": "变量名称",
        "type": "inputString",
        "required": true
      },
      {
        "key": "useInput",
        "name": "是否使用输入值设置全局变量",
        "type": "checkbox",
        "default": true
      },
      {
        "key": "value",
        "name": "变量值",
        "type": "inputString"
      }
    ],
    "ports": {
      "in": [
        {
          "id": "in1",
          "name": "in1"
        }
      ],
      "out": [
        {
          "id": "out1",
          "name": "out1"
        }
      ]
    }
  },
  {
    "type": "statement",
    "typeLabel": "逻辑语句",
    "category": null,
    "categoryLabel": null,
    "name": "获取全局变量",
    "key": "GlobalVarsGetting",
    "helpUrl": null,
    "parameters": [
      {
        "key": "name",
        "name": "变量名称",
        "type": "inputString",
        "required": true
      }
    ],
    "ports": {
      "in": [
        {
          "id": "in1",
          "name": "in1"
        }
      ],
      "out": [
        {
          "id": "out1",
          "name": "out1"
        }
      ]
    }
  },
  {
    "type": "web",
    "typeLabel": "网络协议",
    "category": null,
    "categoryLabel": null,
    "name": "Http请求",
    "key": "HttpRequest",
    "helpUrl": null,
    "parameters": [
      {
        "key": "method",
        "name": "请求方式",
        "options": [
          {
            "value": "post",
            "label": "POST"
          },
          {
            "value": "get",
            "label": "GET"
          }
        ],
        "type": "select",
        "default": "get"
      },
      {
        "key": "address",
        "name": "请求地址",
        "type": "inputString",
        "required": true
      },
      {
        "key": "data",
        "name": "data",
        "type": "inputString",
        "default": ""
      },
      {
        "key": "json",
        "name": "json",
        "type": "inputString",
        "default": ""
      }
    ],
    "ports": {
      "in": [
        {
          "id": "in1",
          "name": "in1"
        }
      ],
      "out": [
        {
          "id": "out1",
          "name": "out1"
        },
        {
          "id": "out2",
          "name": "out2"
        }
      ]
    }
  },
  {
    "type": "triggers",
    "typeLabel": "触发器",
    "category": null,
    "categoryLabel": null,
    "name": "定时器",
    "key": "TriggersComponent",
    "helpUrl": null,
    "parameters": [
      {
        "key": "startDate",
        "name": "开始时间",
        "type": "inputString",
        "required": true
      },
      {
        "key": "endDate",
        "name": "结束时间",
        "type": "inputString"
      },
      {
        "key": "timeUnit",
        "name": "触发时间单位",
        "options": [
          {
            "value": "day",
            "label": "天"
          },
          {
            "value": "hour",
            "label": "小时"
          },
          {
            "value": "minute",
            "label": "分钟"
          },
          {
            "value": "second",
            "label": "秒"
          }
        ],
        "type": "select",
        "default": "second"
      },
      {
        "key": "cronTime",
        "name": "间隔时长",
        "type": "inputInteger",
        "required": true
      }
    ],
    "ports": {
      "in": [
        {
          "id": "in1",
          "name": "in1"
        }
      ],
      "out": [
        {
          "id": "out1",
          "name": "out1"
        }
      ]
    }
  },
  {
    "type": "triggers",
    "typeLabel": "触发器",
    "category": null,
    "categoryLabel": null,
    "name": "等间隔触发定时器",
    "key": "IntervalTriggerComponent",
    "helpUrl": null,
    "parameters": [
      {
        "key": "timeUnit",
        "name": "触发时间单位",
        "options": [
          {
            "value": "days",
            "label": "天"
          },
          {
            "value": "hours",
            "label": "小时"
          },
          {
            "value": "minutes",
            "label": "分钟"
          },
          {
            "value": "seconds",
            "label": "秒"
          }
        ],
        "type": "select",
        "default": "seconds"
      },
      {
        "key": "interval",
        "name": "间隔时长",
        "type": "inputFloat",
        "required": true,
        "default": 5
      }
    ],
    "ports": {
      "out": [
        {
          "id": "out1",
          "name": "out1"
        }
      ]
    }
  },
  {
    "type": "stream",
    "typeLabel": "流计算接口",
    "category": "in",
    "categoryLabel": "输入",
    "name": "输入1",
    "key": "StreamIn_in1",
    "parameters": [],
    "ports": {
      "out": [
        {
          "id": "out1",
          "name": "out1"
        }
      ]
    }
  },
  {
    "type": "stream",
    "typeLabel": "流计算接口",
    "category": "in",
    "categoryLabel": "输入",
    "name": "in2",
    "key": "StreamIn_in2",
    "parameters": [],
    "ports": {
      "out": [
        {
          "id": "out1",
          "name": "out1"
        }
      ]
    }
  },
  {
    "type": "stream",
    "typeLabel": "流计算接口",
    "category": "in",
    "categoryLabel": "输入",
    "name": "in3",
    "key": "StreamIn_in3",
    "parameters": [],
    "ports": {
      "out": [
        {
          "id": "out1",
          "name": "out1"
        }
      ]
    }
  },
  {
    "type": "stream",
    "typeLabel": "流输入与输出",
    "category": "out",
    "categoryLabel": "输出",
    "name": "输出1",
    "key": "StreamOut_out1",
    "parameters": [],
    "ports": {
      "in": [
        {
          "id": "in1",
          "name": "in1"
        }
      ]
    }
  },
  {
    "type": "stream",
    "typeLabel": "流输入与输出",
    "category": "out",
    "categoryLabel": "输出",
    "name": "out2",
    "key": "StreamOut_out2",
    "parameters": [],
    "ports": {
      "in": [
        {
          "id": "in1",
          "name": "in1"
        }
      ]
    }
  },
  {
    "type": "stream",
    "typeLabel": "流输入与输出",
    "category": "out",
    "categoryLabel": "输出",
    "name": "out3",
    "key": "StreamOut_out3",
    "parameters": [],
    "ports": {
      "in": [
        {
          "id": "in1",
          "name": "in1"
        }
      ]
    }
  }
]