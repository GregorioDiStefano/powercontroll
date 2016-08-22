print("Running file mqtt")
tmr.delay(1000000)

-- 6 and 7
gpio12 = 7
gpio.mode(gpio12, gpio.OUTPUT)

gpio13 = 7
gpio.mode(gpio13, gpio.OUTPUT)

mac_address = wifi.sta.getmac()

set_device_status = function(status)
  ok, status_json = pcall(cjson.encode, {status=status, mac_address=wifi.sta.getmac()})
  m:publish("status_alert",status_json, 2, 1, function(client)
    print ("Status alert sent: " .. status)
  end)
end

-- init mqtt client with keepalive timer 120sec
m = mqtt.Client("nodemcu", 120, "@@MQTT_USERNAME@@", "@@MQTT_PASSWORD@@")

ok, json = pcall(cjson.encode, {mac_address=wifi.sta.getmac()})

-- setup Last Will and Testament (optional)
-- Broker will publish a message with qos = 0, retain = 0, data = "offline"
-- to topic "/lwt" if client don't send keepalive packet
m:lwt("/lwt", "offline", 0, 0)

m:on("connect", function(client)
    print ("Connected to MQTT broker.")
end)

m:on("offline", function(client)
  print "Offline.. try reconnecting"
end)


-- on publish message receive event
m:on("message", function(conn, topic, data)
  print(topic .. ":" )

  if topic == "ping" then
    print("Incoming ping")
    ok, pong_json = pcall(cjson.encode, {mac_address=wifi.sta.getmac()})
    m:publish("pong", pong_json, 0,0, function(client)
      print("Pong sent to broker" .. pong_json)
    end)
  end

  if topic == "device/" .. mac_address then
    if pcall(cjson.decode, data) then
      t = cjson.decode(data)
      if t["power"] == "true" then
          print("Turning relay on.")
          set_device_status("on")
          gpio.write(gpio12, gpio.HIGH)
      elseif t["power"] == "false" then
          print("Turning relay off.")
          set_device_status("off")
          gpio.write(gpio12,gpio.LOW)
      end
    else
      print("Error parsing JSON: " .. data)
    end
  end
end)


m:connect("@@MQTT_SERVER_HOSTNAME@@", @@MQTT_PORT@@, 0, 1, function(conn, failed)

    if failed ~= nil then
        print("Connected failed.")
        dofile("init.lua")
        return
    end

    m:subscribe({["device/" .. mac_address ]=0, ["ping"]=0}, function(conn)
        print("Subscribe relay successful.")
    end)

    m:publish("announce",json, 0, 0, function(client) print("Announced device to broker") end)
end)
