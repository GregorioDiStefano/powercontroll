wifi.setmode(wifi.SOFTAP)
wifi.ap.config({ssid="Setup Mode"});
ap_payload=[[
@@AP_PAYLOAD_PAGE@@
]]

function starts(String,Start)
   return string.sub(String,1,string.len(Start))==Start
end

function trim(s)
  return (s:gsub("^%s*(.-)%s*$", "%1"))
end

local function unescape(str)
    str = string.gsub( str, '%%20', ' ')
    str = string.gsub( str, '%%21', '!')
    str = string.gsub( str, '%%22', '"')
    str = string.gsub( str, '%%23', '#')
    str = string.gsub( str, '%%24', '$')
    str = string.gsub( str, '%%25', '%')
    str = string.gsub( str, '%%26', '&')
    str = string.gsub( str, '%%27', '\'')
    str = string.gsub( str, '%%28', '(')
    str = string.gsub( str, '%%29', ')')
    str = string.gsub( str, '%%2A', '*')
    str = string.gsub( str, '%%2B', '+')
    str = string.gsub( str, '%%2C', ',')
    str = string.gsub( str, '%%2D', '-')
    str = string.gsub( str, '%%2E', '.')
    str = string.gsub( str, '%%2F', '.')
    return str
end

srv=net.createServer(net.TCP)
srv:listen(80,function(conn)
    credentials = {}

    conn:on("receive", function(client,request)
      print(request)
      if starts(request, "POST") then
        start = request.find(request, '\r\n\r\n')
        parameter_payload = string.sub(request, start + 2, string.len(request))
        _, _, key, value = string.find(parameter_payload, "(.+)=(.+)&")
        credentials[trim(key)] = unescape(value)

        seperator_location = string.find(parameter_payload, "&")
        second_pair = string.sub(parameter_payload, seperator_location + 1, string.len(request))
        _, _, key, value = string.find(second_pair, "(.+)=(.+)")
        credentials[trim(key)] = unescape(value)

        for key,value in pairs(credentials) do
          print("xxx".. key ..  "xxx".. value .."xxx")
        end

        file.open("ap", "w")
        file.writeline(credentials.ssid)
        file.writeline(credentials.passwd)
        file.close()

        client:send("HTTP/1.1 200 OK")
        client:close()

        dofile("init.lua")
      elseif starts(request, "GET") then
        client:send("HTTP/1.0 200 OK\r\nContent-Type: text/html\r\n\r\n"..ap_payload)
        client:close()
      end
    end)
end)
