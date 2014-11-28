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

    public static boolean/*String*/ sendAuthenticationPostMsg(RequestPackage requestPackage) throws Exception {

        final String LOG_TAG = UserFunctions.class.getSimpleName();
        String uri = requestPackage.getUri();
        String email = requestPackage.getParam("email");
        String password = requestPackage.getParam("password");
        String registration_id = requestPackage.getParam("registration_id");

        URL url = new URL(uri);
        HttpURLConnection con = (HttpURLConnection) url.openConnection();
        con.setRequestMethod("POST");

        String urlParameters = "email=" + email + "&password=" + password + "&registration_id=" + registration_id;

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

        //TODO: remove log
        Log.v(LOG_TAG, "cookie: " + cookie);

        if (cookie!= null) {
            Database.database.put("cookie", cookie);
            return true;
        }
        else {
            return false;
        }
        //return cookie;
    }

    private static String getCookie(String cookieName, HttpURLConnection connection) {
        List<String> setCookieList = connection.getHeaderFields().get("Set-Cookie");
        return setCookieList.get(0).substring(0,cookieName.length()+161);
    }

    public static String getSelectedListItemFunctionality(int position) throws Exception {
        String uri;
        switch (position){
            case 0:
                uri = "http://104.131.126.89/getmessages";
                break;
            case 1:
                uri = "http://104.131.126.89/getmessages";
                break;
            case 2:
                uri = "http://104.131.126.89/groupslist";
                break;
            default:
                uri = "http://104.131.126.89/getmessages";
                break;
        }
        Log.v("checkURI", uri);
        String content = sendGetMsg(uri);
        return content;
    }

    public static String sendGetMsg(String uri) throws Exception {
        URL url = new URL(uri);
        String cookie = Database.database.get("cookie");
        Log.v("checkcookie", cookie);

        HttpURLConnection newConnection = (HttpURLConnection) url.openConnection();

        Log.v("check1", "1");

        newConnection.setRequestMethod("GET");

        newConnection.setRequestProperty("Cookie", cookie);

        Log.v("check2", "2");
        BufferedReader in = new BufferedReader(new InputStreamReader(newConnection.getInputStream()));
        String inputLine;
        StringBuffer response = new StringBuffer();

        while ((inputLine = in.readLine()) != null) {
            response.append(inputLine);
        }
        in.close();

        Log.v("response", response.toString());
        return response.toString();
    }
}


