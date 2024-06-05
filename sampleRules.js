const req = {
    cardNumber: "201221",
    initialBalanace: "100",
}

const rules = [[
    {
        "id": "uuid",
        "op1": "abc",
        "opnd": "eq",
        "op2": "abc",
        "then": [{
            "type": "rule",
            "data": "uuid"
        }],
        "else": [{
            "type": "rule",
            "data": "2"
        }]
    },
]]