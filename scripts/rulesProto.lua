local cjson = require "cjson"

local operators = {
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
local allowedRequestMethods = {"get", "post"}

local function tableIncludes(tbl, value, key)
    for k, v in pairs(tbl) do
        if value ~= nil and v == value then
            return true
        elseif key ~= nil and k == key then
            return true
        end
    end
    return false
end

local jsonRule = [[
    {
        "op1": {
            "name": "{field1}",
            "position": [
                {
                    "op1": "{field3}",
                    "operator": "eq",
                    "op2": {
                        "name": "{field5}",
                        "position": [
                            3,
                            {
                                "op1": "21",
                                "operator": "eq",
                                "op2": "21",
                                "then": "{field6}",
                                "else": 23
                            }
                        ]
                    },
                    "then": "{field4}",
                    "else": 9
                },
                9
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
        "then": 20,
        "else": 200,
        "elseAction": {
            "type": "rest else",
            "msg": {
                "method": "post",
                "url": "localhost:3000",
                "headers": {},
                "data": {}
            }
        }
    }
]]

local sampleReq = [[
    {
        "field1": "abcdefg19klmn",
        "field2": "10",
        "field3": "pss",
        "field4": "8",
        "field5": "abpssde",
        "field6": "5"
    }
]]

req = cjson.decode(sampleReq)

local function validate(rule)
    if type(rule) == "table" then

        if not tableIncludes(rule, nil, "op1") or not tableIncludes(rule, nil, "operator") or
            not tableIncludes(rule, nil, "op2") or
            (not tableIncludes(rule, nil, "then") and not tableIncludes(rule, nil, "thenAction")) or
            (tableIncludes(rule, nil, "then") and tableIncludes(rule, nil, "thenAction")) then
            return false;
        end

        for key, value in pairs(rule) do
            if key == "op1" or key == "op2" then
                if type(value) == "table" then
                    if not value.name or type(value.name) ~= "string" then
                        return false
                    end
                    if type(value.position[1]) == "table" then
                        return validate(value.position[1])
                    elseif type(value.position[1]) ~= "number" then
                        return false
                    end
                    if type(value.position[2]) == "table" then
                        return validate(value.position[2])
                    elseif type(value.position[2]) ~= "number" then
                        return false
                    end
                end
            elseif key == "operator" then
                return tableIncludes(operators, value);
            elseif key == "thenAction" then
                if type(value) ~= "table" then
                    return false
                elseif type(value.type) ~= "string" then
                    return false
                elseif type(value.msg) ~= "table" then
                    return false
                elseif type(value.msg.method) ~= "string" or not tableIncludes(allowedRequestMethods, value.msg.method) then
                    return false
                elseif type(value.msg.url) ~= "string" then
                    return false
                elseif type(value.msg.headers) ~= "table" then
                    return false
                elseif type(value.msg.data) ~= "table" then
                    return false
                end
            end
        end
        return true
    end
    return false
end

local function handleAction(action)
end

local function returnMethod(returnField, action)
    if action ~= nil then
        handleAction(action)
    end
    if type(returnField) == "string" and returnField:match("^%b{}$") ~= nil then
        return req[returnField:match("{(.-)}")]
    else
        return returnField
    end
end

local function evaluate(rule)
    local ops = {rule.op1, rule.op2};
    local opRes = {rule.op1, rule.op2};
    local operator = rule.operator;

    for i, value in pairs(ops) do
        if type(value) == "table" then
            local fieldName = value.name:match("{(.-)}");
            local positions = {0, 0};
            for j, pos in pairs(value.position) do
                if type(pos) == "table" then
                    positions[j] = evaluate(pos)
                else
                    positions[j] = pos
                end
            end
            opRes[i] = req[fieldName]:sub(positions[1], positions[2])
        else
            if type(value) == "string" and value:match("^%b{}$") ~= nil then
                opRes[i] = req[value:match("{(.-)}")]
            end
        end
    end
    if operators[operator](opRes[1], opRes[2]) then
        return returnMethod(rule["then"], rule.thenAction)
    else
        return returnMethod(rule["else"], rule.elseAction)
    end
end

local function main()
    local rule = cjson.decode(jsonRule)
    local start = os.clock()
    print(validate(rule));
    print(evaluate(rule))
    print('Time taken: '.. (os.clock() - start) ..' seconds.')
end

main();
