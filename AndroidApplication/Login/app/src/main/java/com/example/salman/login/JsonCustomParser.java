package com.example.salman.login;

import org.json.JSONObject;

/**
 * Created by salman on 11/24/14.
 */
public class JsonCustomParser {
    public void readAndParseJSON(String in) {
        try {
            JSONObject reader = new JSONObject(in);

            JSONObject sys = reader.getJSONObject("sys");
            //country = sys.getString("country");

            JSONObject main = reader.getJSONObject("main");
            //temperature = main.getString("temp");

            // pressure = main.getString("pressure");
            //humidity = main.getString("humidity");

            // parsingComplete = false;


        } catch (Exception e) {
            e.printStackTrace();
        }
    }
}