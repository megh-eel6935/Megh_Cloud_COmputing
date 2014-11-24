package com.example.salman.login;

import java.io.BufferedReader;
import java.io.DataOutputStream;
import java.io.InputStreamReader;
import java.net.HttpURLConnection;
import java.net.URL;


/**
 * Created by salman on 11/23/14.
 */
public class UserFunctions {

    public static String sendPost(String uri, String email, String password) throws Exception {

        final String USER_AGENT = "Mozilla/5.0";

        URL url = new URL(uri);
        HttpURLConnection con = (HttpURLConnection) url.openConnection();

        //add request header
        con.setRequestMethod("POST");
        //con.setRequestProperty("User-Agent", USER_AGENT);
        //con.setRequestProperty("Accept-Language", "en-US,en;q=0.5");

        String urlParameters = "email=" + email + "&password=" + password;



        // Send post request
        con.setDoOutput(true);
        DataOutputStream wr = new DataOutputStream(con.getOutputStream());
        wr.writeBytes(urlParameters);
        wr.flush();
        wr.close();

        int responseCode = con.getResponseCode();

        BufferedReader in = new BufferedReader(new InputStreamReader(con.getInputStream()));
        String inputLine;
        StringBuffer response = new StringBuffer();

        while ((inputLine = in.readLine()) != null) {
            response.append(inputLine);
        }
        in.close();



        return response.toString();
    }
}

