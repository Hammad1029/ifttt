const req = {
    cardNumber: "201221",
    initialBalanace: "100",
}

const createCard = {
    name: "createCard",
    type: "pre",
    rule: {
        op1: "true",
        operator: "eq",
        op2: "true",
        thenActions: [{
            type: "modifyReq",
            data: {
                field: "currentBalance",
                value: {
                    get: "initialBalance",
                    from: "req"
                }
            }
        }]
    }
}

const addFunds = {
    name: "addFunds",
    type: "pre",
    rule: {
        op1: "true",
        operator: "eq",
        op2: "true",
        thenAction: [{
            type: "updateDb",
            data: {
                table: "createCard",
                where: {
                    cardNumber: {
                        get: "cardNumber",
                        from: "req"
                    }
                },
                set: {
                    field: "balance",
                    type: "add",
                    value: {
                        get: "amount",
                        from: "req"
                    }
                }
            }
        }]
    }
}

const checkBalance = {
    name: "checkBalance",
    type: "pre",
    rule: {
        op1: {
            get: "tnxAmount",
            from: "req"
        },
        operator: "lte",
        op2: {
            get: "balance",
            from: "db",
            condition: {
                table: "createCard",
                where: {
                    cardNumber: {
                        get: "cardNumber",
                        from: "req"
                    }
                }
            }
        },
        thenAction: [{
            type: "update",
            data: {
                table: "createCard",
                where: {
                    cardNumber: {
                        get: "cardNumber",
                        from: "req"
                    }
                },
                set: {
                    field: "balance",
                    type: "subtract",
                    value: {
                        get: "tnxAmount",
                        from: "req"
                    }
                }
            }
        }],
        elseAction: [{
            type: "email",
            data: {
                subject: "Insufficient Balance",
                body: "Pese daal"
            }
        }],
        then: "Transaction Successful",
        else: "Insufficient Balance"
    }
}

const signUp = {
    name: "signup",
    type: "pre",
    rule: {
        op1: {
            get: "email",
            from: "db",
            condition: {
                query: {
                    email: {
                        get: "email",
                        from: "req",
                    }
                },
                apiName: "signup"
            }
        },
        operator: "eq",
        op2: "nil",
        thenActions: [
            {
                type: "insert",
                data: {
                    get: "*",
                    from: "req",
                    condition: {}
                }
            },
            {
                type: "insert",
                data: {
                    get: "province",
                    from: "rule",
                    condition: {
                        rule: "form_of_identity"
                    }
                }
            },
            {
                type: "email",
                data: {
                    subject: "Signup verification",
                    body: "You have been signed up"
                }
            }
        ],
        elseActions: [],
        then: "true",
        else: "false"
    }
}

const province = {
    name: "form_of_identity",
    type: "pre",
    rule: {
        op1: {
            get: "cnic",
            from: "req",
            condition: {
                mutate: ["substring(0,1)"]
            }
        },
        operator: "eq",
        op2: "4",
        thenAction: [],
        elseAction: [],
        then: "Sindh",
        else: {
            get: "city",
            from: "req",
        }
    }
}

const checkPassword = {
    name: "check_user_password",
    type: "pre",
    rule: {
        op1: {
            get: {
                field: "password",
                from: "req",
                condition: {
                    mutate: ["md5"]
                }
            }
        },
        operand: "eq",
        op2: {
            get: "password",
            from: "db",
            condition: {
                query: {
                    email: {
                        get: "email",
                        from: "req"
                    }
                }
            }
        }
    }
}

const oldRule = {
    "rule": {
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
    },
    "req": {
        "field1": "abcdefg19klmn",
        "field2": "10",
        "field3": "pss",
        "field4": "8",
        "field5": "abpssde",
        "field6": "5"
    }
}