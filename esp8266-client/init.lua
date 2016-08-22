wifi.setmode(wifi.STATION)
wifi.sta.disconnect()

failedConnectionAttempts = 0

AP_SSID = nil
AP_PASSWORD = nil


function trim(s)
  return (s:gsub("^%s*(.-)%s*$", "%1"))
end

if file.exists("ap") then
    file.open("ap", "r")
    AP_SSID = trim(file.readline())
    AP_PASSWORD = trim(file.readline())
    print("Loaded, ssid: " ..AP_SSID .." passwd:" .. AP_PASSWORD)
    file.close()
end

tmr.alarm(1, 10000, 1, function()
    print("Checking connection, attempt: "..failedConnectionAttempts)
    if wifi.sta.getip()== nil or internet_works()== nil then
        if failedConnectionAttempts >= 15 then
          print("Connection failed too many times, enabling AP mode.")
          dofile("apsetup.lua")
          tmr.stop(1)
          return
        end

        if file.exists("ap") == false then
          print("No ap file found, starting AP")
          dofile("apsetup.lua")
          tmr.stop(1)
          return
        end

        print "No IP or DNS lookup failed, retry connecting to WiFi"
	      connect()
    else
        tmr.stop(1)
        failedConnectionAttempts = 0
        print("Config done, IP is "..wifi.sta.getip())
        dofile("mqtt2.lua")
        return
    end
 end)


internet_works = function ()
	local retVal = 0
	net.dns.setdnsserver("8.8.8.8", 0)
	net.dns.setdnsserver("8.8.4.4", 1)
	net.dns.resolve("www.google.com", function(sk, ip)
  if (ip == nil) then
		     retVal = nil
	   else
		     retVal = 1
	   end
	end)
	return retVal
end

connect = function ()
	--- move this to a on "connected" callback
	print ("Trying to connect to AP.........")
	wifi.setmode(wifi.STATION)
	wifi.sta.config(AP_SSID, AP_PASSWORD, 1)
	wifi.sta.connect()
end

wifi.sta.eventMonReg(wifi.STA_WRONGPWD, function()
  print("STATION_WRONG_PASSWORD")
  failedConnectionAttempts = failedConnectionAttempts + 1
  connect()
  tmr.start(1)
end)

wifi.sta.eventMonReg(wifi.STA_APNOTFOUND, function()
  print("STATION_NO_AP_FOUND")
  failedConnectionAttempts = failedConnectionAttempts + 1
  connect()
  tmr.start(1)
end)

wifi.sta.eventMonReg(wifi.STA_FAIL, function()
  print("STATION_CONNECT_FAIL")
  failedConnectionAttempts = failedConnectionAttempts + 1
  connect()
  tmr.start(1)
end)

wifi.sta.eventMonStart(100)
