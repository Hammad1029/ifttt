local utils = {}

function utils.printTable(table)
    for k, v in pairs(table) do
        print(k, v)
    end
    print()
end

function utils.getValue(getSchema)
    if getSchema.from == "req" then
        return req[getSchema.get]
    end
end

function utils.invertMapping(mapping)
    local s = {}
    for k, v in pairs(mapping) do
        s[v] = k
    end
    return s
end

function utils.mapReqToInternal(mapping)
    local mapped = {}
    for key, value in pairs(req) do
        local internalCol = mapping[key]
        if internalCol ~= nil then
            mapped[internalCol] = value
        else
            mapped[key] = value
        end
    end
    return mapped
end

return utils