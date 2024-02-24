local json = require("./scripts/jsonParse")
local query = require("./scripts/query")
local utils = require("./scripts/utils")

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
    return comparisionOperators[operator](ops[1], ops[2])
end

local function handleActions(actions)
    for _, ac in pairs(actions) do
        if ac.type == "modifyReq" then
            req[ac.data.field] = utils.getValue(ac.data.value)
        elseif ac.type == "updateDb" then
            query.update(ac.data)
        elseif ac.type == "insertDb" then
            query.insert(ac.data)
        end
    end
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

function evaluate(rule, req, mapping)
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

local function parseAndEvaluate()
    return reqJson
    -- local rule = json.parse(ruleJson).rule
    -- req = json.parse(reqJson)
    -- local mapping = json.parse(apiMapping)
    -- evaluate(rule, apiMapping)
    -- return json.stringify(utils.mapReqToInternal(mapping))
end

req = {}

print(parseAndEvaluate())