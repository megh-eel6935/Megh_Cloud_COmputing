package com.example.salman.login;

import android.util.Log;

import java.io.BufferedReader;
import java.io.DataOutputStream;
import java.io.InputStreamReader;
import java.net.HttpURLConnection;
import java.net.URL;
import java.util.List;


/**
 * Created by Ahmad Salman Saqib.
 */
public class UserFunctions {

    public static String sendAuthenticationPost(RequestPackage requestPackage) throws Exception {

        final String LOG_TAG = UserFunctions.class.getSimpleName();
        String uri = requestPackage.getUri();
        String email = requestPackage.getParam("email");
        String password = requestPackage.getParam("password");

        URL url = new URL(uri);
        HttpURLConnection con = (HttpURLConnection) url.openConnection();
        con.setRequestMethod("POST");

        String urlParameters = "email=" + email + "&password=" + password;

        con.setDoOutput(true);
        DataOutputStream wr = new DataOutputStream(con.getOutputStream());
        wr.writeBytes(urlParameters);
        wr.flush();
        wr.close();

        BufferedReader in = new BufferedReader(new InputStreamReader(con.getInputStream()));
        String inputLine;
        StringBuffer response = new StringBuffer();

        while ((inputLine = in.readLine()) != null) {
            response.append(inputLine);
        }
        in.close();

        String cookie = getCookie("session", con);

        Log.v(LOG_TAG, "cookie: " + cookie);

        return cookie;
    }


    public static String sendGet(String uri, String cookie) throws Exception {
        URL urlGroup = new URL(uri);
        HttpURLConnection newConnection = (HttpURLConnection) urlGroup.openConnection();
        //add request header
        newConnection.setRequestMethod("GET");

        newConnection.setRequestProperty("Cookie", cookie);

        BufferedReader inGroup = new BufferedReader(new InputStreamReader(newConnection.getInputStream()));
        String inputLineGroup;
        StringBuffer responseGroup = new StringBuffer();

        while ((inputLineGroup = inGroup.readLine()) != null) {
            responseGroup.append(inputLineGroup);
            Log.v(UserFunctions.class.getSimpleName(), inputLineGroup);
        }
        inGroup.close();

        return responseGroup.toString();
    }

    private static String getCookie(String cookieName, HttpURLConnection connection) {
        List<String> setCookieList = connection.getHeaderFields().get("Set-Cookie");
        return setCookieList.get(0).substring(0,cookieName.length()+161);
    }
}


