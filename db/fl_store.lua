    require'strict'.on()

box.cfg{
    listen = 3000
}

box.schema.user.grant('guest', 'read,write,execute', 'universe', nil,{
    if_not_exists=true
})

--Создание спейса юзера
user = box.schema.create_space('user', {
    if_not_exists = true,
})

user:format(
        {
            {name = 'id', type = 'unsigned'},
            {name = 'email', type = 'string'},
            {name = 'password', type = 'string'},
            {name = 'user_name', type = 'string'},
            {name = 'first_name', type = 'string'},
            {name = 'second_name', type = 'string'},
            {name = 'executor', type = 'boolean'},
            {name = 'description', type = 'string', is_nullable=true},
            {name = 'specializes', type = 'array', is_nullable=true},
            {name = 'img_url', type = 'string', is_nullable = true},
        }
)

user:create_index('primary', {
    sequence = true,
    type = 'TREE',
    parts = {'id'},
    if_not_exists = true,
})

user:create_index('email_key', {
    unique = true,
    type = 'HASH',
    parts = {'email'},
    if_not_exists = true,
})

--Создание спейса заказов
order = box.schema.create_space('order', {
    if_not_exists = true
})

order:format(
        {
            {name = 'id',          type = 'unsigned'},
            {name = 'order_name',  type = 'string'},
            {name = 'customer_id', type = 'unsigned'},
            {name = 'description', type = 'string'},
            {name = 'specializes', type = 'array' },
        }
)

order:create_index('primary', {
    sequence = true,
    type = 'TREE',
    parts = {'id'},
    if_not_exists = true,
})

order:create_index('customer_index', {
    type = 'TREE',
    parts = {'customer_id'},
    if_not_exists = true,
    unique = false,
})

order:create_index('order_name', {
    type = 'TREE',
    parts = {'order_name'},
    if_not_exists = true,
    unique = false,
})



--Создание спейса специализаций
specialize = box.schema.create_space('specialize', {
    if_not_exists = true
})

specialize:format(
        {
            {name = 'specialize_name', type = 'string'},
            {name = 'customer_list', type = 'array'},
            {name = 'order_list', type = 'array'},
        }
)

specialize:create_index('primary', {
    type = 'HASH',
    parts = {'specialize_name'},
    if_not_exists = true,
})

--Хранилище кук
session = box.schema.create_space('session', {
    if_not_exists = true
})

session:format(
        {
            {name = 'session', type = 'string'},
            {name = 'user_id', type = 'unsigned'},
        }
)

session:create_index('primary', {
    type = 'HASH',
    parts = {'session'},
    if_not_exists = true,
})

function alldrop()
    user:drop()
    order:drop()
    specialize:drop()
    session:drop()
end

function GetOrders(num, spec)
    local count = 0;
    for _, v in pairs(spec) do
        spec[v] = true
    end
    local result = {}
    for _, v in order:pairs() do
        if count == num then
            break
        end
        local specialize = v[5]
        for _, elem in pairs(specialize) do
            local flag = spec[elem]
            if flag then
                table.insert(result, count, v)
                break
            end
        end
        count = count + 1
    end
    return box.tuple.new(result)
end

function Test()
    local testTuple = box.tuple.new()
    local count = 0

    --for _, v in order:pairs() do
    --    testTuple[count] = v
    --    count = count + 1
    --end

    return testTuple
end

require'console'.start()