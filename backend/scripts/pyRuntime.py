import json
import argparse
import pandas
import requests
import traceback

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

loadMethods = {
    "string": str,
    "int": int,
    "float": float,
    "json": defaultLoad,
    "bool": defaultLoad
}

def defaultDump(x):
    return x

dumpMethods = {
    str: defaultDump,
    int: defaultDump,
    float: defaultDump,
    dict: defaultDump,
    list: defaultDump,
    bool: defaultDump
}

typeMappings = {
    str: "json",
    int: "json",
    float: "json",
    dict: "json",
    list: "json",
    bool: "json"
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

def run(inputs=None, script=""):
    exec(functionSting % script.replace("\n", "\n    "), globals())
    loadedInputs = []
    for input in inputs:
        input = json.loads(eval("'{}'".format(input)))
        loadedInputs.append(loadMethods[input["type"]](input["data"]))
    outputs = runScript(loadedInputs)
    dumpedOutputs = []
    for output in outputs:
        if type(output) in dumpMethods:
            dumpedOutputs.append({"data": dumpMethods[type(output)](output), "type": typeMappings[type(output)]})
        elif not output:
            dumpedOutputs.append({"data": output, "type": "json"})
        else:
            raise Exception(f"type of {output} is not supported.")
    return json.dumps(dumpedOutputs)

def run(inputs=None, script=""):
    exec(functionSting % script.replace("\n", "\n    "), globals())
    loadedInputs = []

    for input in inputs:
        # print(input)
        input = json.loads(eval("'{}'".format(input)))
        loadedInputs.append(loadMethods[input["type"]](input["data"]))
    # input = json.loads(eval("'{}'".format(inputs)))
    # loadedInputs.append(loadMethods[input["type"]](input["data"]))
    outputs = runScript(loadedInputs)
    dumpedOutputs = []
    for output in outputs:
        if type(output) in dumpMethods:
            dumpedOutputs.append({"data": dumpMethods[type(output)](output), "type": typeMappings[type(output)]})
        else:
            raise Exception(f"type of {output} is not supported.")
    return dumpedOutputs


@app.get("/data/")
async def getInputdata(inputdata, script):
    tmp = inputdata.split("},")
    if len(tmp) > 1:
        for i in range(len(tmp) - 1):
            tmp[i] = tmp[i] + "}"
    result = run(tmp, script)
    print(result)
    return result

if __name__=="__main__":
    uvicorn.run(app, host="0.0.0.0", port=8080)
