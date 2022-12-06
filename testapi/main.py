from fastapi import FastAPI

app = FastAPI()
# inputdata = {"data":"1","type":"json"}
# script = "# 输入端口数据通过inputs[n]来引用n从0开始:\n# 输处端口数据放入outputs，outputs为列表，按端口顺序放入数据\n# outputs.append(int(inputs[0])+1)\noutputs.append(int(inputs[0])+2)"
# nodeid = ""
import json
import argparse
# import pandas


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

def run(inputs=None, script=""):
    exec(functionSting % script.replace("\n", "\n    "), globals())
    loadedInputs = []
    print("ly---")

    for input in inputs:
        print(input)
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
    return json.dumps(dumpedOutputs)


@app.get("/data/")
async def getInputdata(inputdata, script):
    print(inputdata)#{"data":"1","type":"json"},{"data":"2","type":"json"}
    print(script)
    tmp = inputdata.split("},")
    tmp[0] = tmp[0] + "}"
    print(tmp)
    result = run(tmp, script)
    # return {
    #     "input": inputdata,
    #     "script": script
    # }

    return result
