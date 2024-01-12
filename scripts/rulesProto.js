

const sampleData = {
    rule: JSON.stringify({
        "op1": {
            "name": "{field1}",
            "position": [
                {
                    "op1": "{field3}",
                    "operator": "eq",
                    "op2": {
                        "name": "{field5}",
                        "position": [
                            "3",
                            {
                                "op1": "21",
                                "operator": "eq",
                                "op2": "21",
                                "then": "{field6}",
                                "else": "23"
                            }
                        ]
                    },
                    "then": "{field4}",
                    "else": "9"
                },
                "9"
            ]
        },
        "operator": "gte",
        "op2": "{field2}",
        "thenAction": {
            "type": "rest",
            "msg": {
                "method": "post",
                "url": "localhost:3000",
                "headers": {},
                "data": {}
            }
        },
        "then": "20",
        "else": "200",
        "elseAction": {
            "type": "rest else",
            "msg": {
                "method": "post",
                "url": "localhost:3000",
                "headers": {},
                "data": {}
            }
        }
    }),
    req: JSON.stringify({
        "field1": "abcdefg19klmn",
        "field2": "10",
        "field3": "pss",
        "field4": "8",
        "field5": "abpssde",
        "field6": "5"
    })
}

const schemas = {
    rule: {
        op1: ["string", "object"],
        op2: ["string", "object"],
        operator: ["string"],
        then: ["string"],
        thenAction: ["object"],
        else: ["string"],
        elseAction: ["object"],
    },
    subRule: {
        op1: ["string", "object"],
        op2: ["string", "object"],
        operator: ["string"],
        then: ["string"],
        else: ["string"],
    },
    operand: {
        name: ["string"],
        position: ["array", 2, ["string", "object"]]
    },
    actions: {
        type: ["string"],
        msg: ["object"],
    },
    actionMsg: {
        method: ["string"],
        url: ["string"],
        headers: ["object"],
        data: ["object"]
    }
}

const validateSchemas = (obj, schema) => Object.keys(schema).every(key =>
    Array.isArray(obj[key])
        ? schema[key].includes("array") && obj[key].length === schema[key][1]
        && obj[key].every(i => schema[key][2].includes(typeof i))
        : schema[key].includes(typeof obj[key])
)

const operators = ([x = 0, y = 0] = [0, 0]) => ({
    lt: () => Number(x) < Number(y),
    gt: () => Number(x) > Number(y),
    gte: () => Number(x) >= Number(y),
    lte: () => Number(x) <= Number(y),
    eq: () => x == y,
    neq: () => x != y
})

const validate = (rule, type = "subRule") => {
    try {
        if (typeof rule === "object") {
            if (validateSchemas(rule, schemas[type])) {
                const valid = Object.keys(rule).every(key => {
                    switch (key) {
                        case "op1":
                        case "op2":
                            const opType = typeof rule[key]
                            if (opType === "object") {
                                if (validateSchemas(rule[key], schemas.operand)) {
                                    return (
                                        typeof rule[key].position[0] === "object" && typeof rule[key].position[1] === "object"
                                            ? validate(rule[key].position[0]) && validate(rule[key].position[1])
                                            : typeof rule[key].position[0] === "object" && typeof rule[key].position[1] !== "object"
                                                ? validate(rule[key].position[0])
                                                : typeof rule[key].position[0] !== "object" && typeof rule[key].position[1] === "object"
                                                    ? validate(rule[key].position[1])
                                                    : false
                                    )
                                }
                            } else if (opType === "string") return true
                            return false
                        case "operator":
                            return Object.keys(operators()).includes(rule[key])
                        case "then":
                        case "else":
                            return true
                        case "thenAction":
                        case "elseAction":
                            return validateSchemas(rule[key], schemas.actions) && validateSchemas(rule[key].msg, schemas.actionMsg)
                        default:
                            return false;
                    }
                })
                return valid
            }
        }
        return false
    } catch (e) {
        console.error(e)
        return false
    }
}

const getField = (str, obj) => {
    const field = String(str).match(/{([^}]+)}/)
    if (field === null) return false
    else return obj[field[1]]
}

const handleAction = (action) => {
    console.log(action)
}

const evaluate = (rule, req) => {
    try {
        const ops = [rule.op1, rule.op2]
        const opRes = ops.map(o => {
            switch (typeof o) {
                case "object":
                    const positions = [o.position[0], o.position[1]]
                    const posRes = positions.map(p => {
                        switch (typeof p) {
                            case "object":
                                return Number(evaluate(p, req))
                            case "string":
                                return Number(p)
                        }
                    })
                    return getField(o.name, req)?.substring(posRes[0] - 1, posRes[1]) || o.name
                case "string":
                    return getField(o, req) || o
                default:
                    throw Error("Wrong type of operand")
            }
        })
        const doAfter = operators(opRes)[rule.operator]() ? "then" : "else"
        const action = rule[doAfter + "Action"]
        action && handleAction(action)
        const returnVal = String(rule[doAfter])
        return getField(returnVal, req) || returnVal
    } catch (e) {
        console.error(e)
        return false
    }
}

const benchmark = (fn) => {
    const start = performance.now();
    const returnVal = fn()
    const end = performance.now();
    const timeTaken = (end - start) / 1000;
    console.log('Time taken: ' + timeTaken + ' seconds.');
    return returnVal
}

const parseAndValidate = (rule) => validate(JSON.parse(rule))
const parseAndEvaluate = (rule, req) => evaluate(JSON.parse(rule), JSON.parse(req))

const main = () => {
    const { rule, req } = sampleData
    console.log(benchmark(() => {
        const validated = parseAndValidate(rule)
        const evaluated = parseAndEvaluate(rule, req)
        return console.log(
            validated,
            evaluated
        )
        // return JSON.stringify({ validated, evalution })
    }))
}

main()