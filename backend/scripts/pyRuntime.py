import json
import argparse

functionSting = '''
def runScript(inputs):
    %s
    return outputs
'''

loadMethods = {
    "string": str,
    "int": int,
    "float": float,
    "json": json.loads
}

dumpMethods = {
    str: str,
    int: str,
    float: str,
    dict: json.dumps,
    list: json.dumps
}

typeMappings = {
    str: "string",
    int: "int",
    float: "float",
    dict: "json",
    list: "json"
}

def run(inputs=None, script=""):
    exec(functionSting % script.replace("\\n", "\n    "), globals())
    loadedInputs = []
    for input in inputs:
        input = json.loads(input)
        loadedInputs.append(loadMethods[input["type"]](input["data"]))
    outputs = runScript(loadedInputs)
    dumpedOutputs = []
    for output in outputs:
        if type(output) in dumpMethods:
            dumpedOutputs.append(json.dumps({"data": dumpMethods[type(output)](output), "type": typeMappings[type(output)]}))
        else:
            raise Exception(f"type of {output} is not supported.")
    return dumpedOutputs

if __name__=="__main__":
    parser = argparse.ArgumentParser(description='Process some user define function.')
    parser.add_argument('inputs', metavar='{"data": 1, "type": "int"}', type=str, nargs='+',
                    help='an input for the script function')
    parser.add_argument('--script', dest='script', type=str,
                        help='script function process the inputs.')
    args = parser.parse_args()
    print(",".join(run(args.inputs, args.script)))
