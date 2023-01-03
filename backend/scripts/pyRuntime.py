import json
import argparse
import pandas as pd
import requests
import traceback
import os

import uvicorn
from fastapi import FastAPI
app = FastAPI()

functionSting = '''
def runScript(inputs):
    outputs = []
    %s
    return outputs
'''
def defaultLoad(x):
    return x

def csvLoad(path):
    return pd.read_csv(path)

loadMethods = {
    "string": str,
    "int": int,
    "float": float,
    "json": defaultLoad,
    "bool": defaultLoad,
    "csv": csvLoad
}

def defaultDump(x):
    return x

def safeMkdirs(path):
    if not os.path.exists(path):
        try:
            os.makedirs(path)
        except FileExistsError:
            pass
    return path

def safeMkdirsForFile(filepath):
    return safeMkdirs(os.path.dirname(os.path.abspath(filepath)))

def csvDump(df, nodeid, idx):
    path = nodeid + "output/"+ "out"+str(idx) +"/data.csv"
    safeMkdirsForFile(path)
    print(df)
    df.to_csv(path, encoding="utf-8", index=True)
    # print(pd.read_csv(path))
    return path

dumpMethods = {
    str: defaultDump,
    int: defaultDump,
    float: defaultDump,
    dict: defaultDump,
    list: defaultDump,
    bool: defaultDump,
    pd.DataFrame: csvDump
}

typeMappings = {
    str: "json",
    int: "json",
    float: "json",
    dict: "json",
    list: "json",
    bool: "json",
    pd.DataFrame: "csv"
}


functionSting = '''
def runScript(inputs):
    outputs = []
    %s
    return outputs
'''
def getGlobalVar(name):
    r = requests.get('http://0.0.0.0:8888/variable', params={"name": name})
    return json.loads(r.content)["data"]

def setGlobalVar(name, data):
    r = requests.post('http://0.0.0.0:8888/variable', params={"name": name}, json=data)
    return json.loads(r.content)

def delGlobalVar(name):
    r = requests.delete('http://0.0.0.0:8888/variable', params={"name": name})
    return json.loads(r.content)

# def run(inputs=None, script=""):
#     exec(functionSting % script.replace("\n", "\n    "), globals())
#     loadedInputs = []
#     for input in inputs:
#         input = json.loads(eval("'{}'".format(input)))
#         loadedInputs.append(loadMethods[input["type"]](input["data"]))
#     outputs = runScript(loadedInputs)
#     dumpedOutputs = []
#     for output in outputs:
#         if type(output) in dumpMethods:
#             dumpedOutputs.append({"data": dumpMethods[type(output)](output), "type": typeMappings[type(output)]})
#         elif not output:
#             dumpedOutputs.append({"data": output, "type": "json"})
#         else:
#             raise Exception(f"type of {output} is not supported.")
#     return json.dumps(dumpedOutputs)

def run(nodeid =None, inputs=None, script=""):
    exec(functionSting % script.replace("\n", "\n    "), globals())
    loadedInputs = []

    for input in inputs:
        # print(input)
        input = json.loads(eval("'{}'".format(input)))
        loadedInputs.append(loadMethods[input["type"]](input["data"]))
    # input = json.loads(eval("'{}'".format(inputs)))
    # loadedInputs.append(loadMethods[input["type"]](input["data"]))
    # print(loadedInputs)
    outputs = runScript(loadedInputs)
    # print(outputs)
    # print(type(outputs))
    dumpedOutputs = []
    idx = 1
    for output in outputs:
        print(type(output))
        print(output)
        if type(output) in dumpMethods:
            if type(output) == pd.DataFrame:
                dumpedOutputs.append({"data": dumpMethods[type(output)](output, nodeid, idx), "type": typeMappings[type(output)]})
                idx += 1
            else:
                dumpedOutputs.append({"data": dumpMethods[type(output)](output), "type": typeMappings[type(output)]})
        elif output is not None:
            dumpedOutputs.append({"data": output, "type": "json"})
        else:
            raise Exception(f"type of {output} is not supported.")
    return dumpedOutputs


@app.get("/data/")
async def getInputdata(nodeid, inputdata, script):
    print(nodeid)
    print(inputdata)
    print(script)

    tmp = inputdata.split("},")
    if len(tmp) > 1:
        for i in range(len(tmp) - 1):
            tmp[i] = tmp[i] + "}"
    result = run(nodeid,tmp, script)
    print(result)
    return result

if __name__=="__main__":
    uvicorn.run(app, host="0.0.0.0", port=8080)
