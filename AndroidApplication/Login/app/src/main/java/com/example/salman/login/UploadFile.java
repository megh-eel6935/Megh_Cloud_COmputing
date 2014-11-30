package com.example.salman.login;

import android.util.Log;

import java.io.BufferedReader;
import java.io.DataOutputStream;
import java.io.File;
import java.io.FileInputStream;
import java.io.InputStreamReader;
import java.net.HttpURLConnection;
import java.net.URL;

public class UploadFile {
    public static boolean handleFile(String filePath) throws Exception {
        HttpURLConnection connection = null;
        DataOutputStream outStream = null;
        InputStreamReader inStream = null;

        String lineEnd = "\r\n";
        String twoHyphens = "--";
        String boundary = "*****";

        int bytesRead, bytesAvailable, bufferSize;

        byte[] buffer;

        int maxBufferSize = 1 * 1024 * 1024;

        String urlString = "http://104.131.126.89/uploadfile";


        FileInputStream fileInputStream = null;

        System.out.println(filePath);

        Log.v("filePath_debug", filePath);

        fileInputStream = new FileInputStream(new File(filePath));

        URL url = new URL(urlString);
        String cookie = Database.database.get("cookie");
        String fileName = Database.database.get("currentFileName");

        connection = (HttpURLConnection) url.openConnection();
        connection.setDoInput(true);
        connection.setDoOutput(true);
        connection.setUseCaches(false);

        connection.setRequestMethod("POST");
        connection.setRequestProperty("Cookie", cookie);
        connection.setRequestProperty("Connection", "Keep-Alive");
        connection.setRequestProperty("Content-Type", "multipart/form-data;boundary=" + boundary);

        outStream = new DataOutputStream(connection.getOutputStream());

        outStream.writeBytes(addParam("file", fileName, twoHyphens, boundary, lineEnd));

        outStream.writeBytes(twoHyphens + boundary + lineEnd);
        outStream.writeBytes("Content-Disposition: form-data; name=\"uploadedfile\"; filename=\"" + filePath + "\"" + lineEnd + "Content-Type: " + "file" + lineEnd + "Content-Transfer-Encoding: binary" + lineEnd);
        outStream.writeBytes(lineEnd);

        bytesAvailable = fileInputStream.available();
        bufferSize = Math.min(bytesAvailable, maxBufferSize);
        buffer = new byte[bufferSize];

        bytesRead = fileInputStream.read(buffer, 0, bufferSize);

        while (bytesRead > 0) {
            outStream.write(buffer, 0, bufferSize);
            bytesAvailable = fileInputStream.available();
            bufferSize = Math.min(bytesAvailable, maxBufferSize);
            bytesRead = fileInputStream.read(buffer, 0, bufferSize);
        }

        outStream.writeBytes(lineEnd);
        outStream.writeBytes(twoHyphens + boundary + twoHyphens + lineEnd);

        fileInputStream.close();
        outStream.flush();
        outStream.close();


        inStream = new InputStreamReader(connection.getInputStream());
        BufferedReader br = new BufferedReader(inStream);
        String str;

        while ((str = br.readLine()) != null) {
            if (JsonCustomParser.readAndParseJSONMessages(str).equals("{\"status\":\"Success\"}")) {
                return true;
            } else {
                return false;
            }
        }
        inStream.close();


        return false;
    }

    private static String addParam(String key, String value, String twoHyphens, String boundary, String lineEnd) {
        return twoHyphens + boundary + lineEnd + "Content-Disposition: form-data; name=\"" + key + "\"" + lineEnd + lineEnd + value + lineEnd;
    }
}