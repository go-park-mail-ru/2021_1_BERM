require 'strict'.on()

box.cfg {
    listen = 3001,
    background = true,
    log = '1.log',
    pid_file = '1.pid'
}

box.schema.user.grant('guest', 'read,write,execute', 'universe', nil, {
    if_not_exists = true
})


--Хранилище кук
session = box.schema.create_space('session', { if_not_exists = true })

session:format(
        {
            { name = 'sessionID', type = 'string' },
            { name = 'user_id', type = 'unsigned' },
            { name = 'executor', type = 'boolean' },
        }
)

session:create_index('primary', {
    type = 'HASH',
    parts = { 'sessionID' },
    if_not_exists = true,
})

function alldrop()
    user:drop()
    order:drop()
    specialize:drop()
    session:drop()
end

require 'console'.start()