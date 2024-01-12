local json = require("jsonParse")

local comparisionOperators = {
    ["lt"] = function(x, y)
        return tonumber(x) < tonumber(y)
    end,
    ["gt"] = function(x, y)
        return tonumber(x) > tonumber(y)
    end,
    ["gte"] = function(x, y)
        return tonumber(x) >= tonumber(y)
    end,
    ["lte"] = function(x, y)
        return tonumber(x) <= tonumber(y)
    end,
    ["eq"] = function(x, y)
        return x == y
    end,
    ["neq"] = function(x, y)
        return x ~= y
    end
}

local function compare(ops, operator)
    return comparisionOperators[operator(ops[1], ops[2])]
end

local function handleActions(actions)
end

local function getFromDb(query)

end

local function getRuleByName(name, req)

end

local function resolveOperand(op, req, mapping)
    if type(op) == "table" then
        if op.from == "req" then
            return req[op.get]
        elseif op.from == "db" then
            return getFromDb(op.condition, mapping)
        elseif op.from == "rule" then
            return evaluate(getRuleByName(op.condition.rule), req)
        end
    else
        return op
    end
end

local function evaluate(rule, req, mapping)
    local opRes = {resolveOperand(rule.op1, req, mapping), resolveOperand(rule.op2, req, mapping)}
    local comparision = compare(opRes, rule.operator)
    if comparision then
        handleActions(rule.thenActions)
        return rule["then"]
    else
        handleActions(rule.elseActions)
        return rule["else"]
    end
end

local function parseAndEvaluate(ruleJson, reqJson, mapping)
    local rule = json.parse(ruleJson).rule
    local req = json.parse(reqJson)
    evaluate(rule, req, mapping)
end

local function main()
    print(parseAndEvaluate(
        [[{"name":"signup","type":"pre","rule":{"op1":{"get":"email","from":"db","condition":{"query":{"email":{"get":"email","from":"req"}}}},"operator":"eq","op2":"nil","thenAction":[{"type":"insert","data":{"get":"*","from":"req","condition":{}}},{"type":"insert","data":{"get":"province","from":"rule","condition":{"rule":"form_of_identity"}}},{"type":"email","data":{"subject":"Signup verification","body":"You have been signed up"}}],"elseAction":[],"then":"true","else":"false"}}]],
        [[{"email":"hammad1029@gmail.com","password":"hello123","name":"Hammad","address":"73 faran","city":"karachi","nationality":"pk","cnic":"4220173029169"}]],
        [[{"name":"String1","age":"int2"}]]))
end

main()
