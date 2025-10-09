-- 你的验证码在 Redis 上的 key
-- email_code:login:user@example.com 替换为实际的邮箱地址
local key = KEYS[1]
-- 验证次数，我们一个验证码，最多重复三次，这个记录还可以验证几次
-- email_code:login:user@example.com:cnt
local cntKey = key..":cnt"
-- 你的验证码 123456
local val = ARGV[1]
-- 过期时间
local ttl = tonumber(redis.call("ttl", key))
if ttl == -1 then
    -- key 存在，但是没有过期时间
    -- 系统错误，你的同事手贱，手动设置了这个 key，但是没给过期时间
    return -2
elseif ttl == -2 or ttl < 540 then
    redis.call("set", key, val)
    redis.call("expire", key, 600) -- 设置600秒的有效期
    redis.call("set", cntKey, 3) -- 设置可以验证3次
    redis.call("expire", cntKey, 600) -- 设置计数器的有效期也是600秒
    -- 完美，符合预期
    return 0
else
    -- 发送太频繁
    return -1
end