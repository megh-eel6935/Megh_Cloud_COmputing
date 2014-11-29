/*
package com.example.salman.login;

import android.util.Log;

import java.io.DataOutputStream;
import java.io.File;
import java.io.FileInputStream;
import java.io.IOException;
import java.io.InputStream;
import java.net.HttpURLConnection;
import java.net.URL;
import java.util.Iterator;
import java.util.Map;

*/
/**
 * Created by salman on 11/28/14.
 *//*

public class UploadFile {

    public ResponseObject multipartRequest(String urlTo, String post, Map httpPara, String filepath, String filefield) throws ParseException, IOException {
        HttpURLConnection connection = null;
        DataOutputStream outputStream = null;
        InputStream inputStream = null;
        int ResponseCode;

        String twoHyphens = "--";
        String boundary =  "*****"+Long.toString(System.currentTimeMillis())+"*****";
        String lineEnd = "\r\n";

        String result = "";

        int bytesRead, bytesAvailable, bufferSize;
        byte[] buffer;
        int maxBufferSize = 1*1024*1024;

        String[] q = filepath.split("/");
        int idx = q.length - 1;

        try {
            File file = new File(filepath);
            FileInputStream fileInputStream = new FileInputStream(file);
            URL url = new URL(urlTo);
            connection = (HttpURLConnection) url.openConnection();
            connection.setDoInput(true);
            connection.setDoOutput(true);
            connection.setUseCaches(false);

            connection.setRequestMethod("POST");
            connection.setRequestProperty("Connection", "Keep-Alive");
            connection.setRequestProperty("User-Agent", "Android Multipart HTTP Client 1.0");
            connection.setRequestProperty("Content-Type", "multipart/form-data; boundary="+boundary);
            Log.i(Commons.TAG, "filepath " + getMimeType(filepath) );

            outputStream = new DataOutputStream(connection.getOutputStream());
            outputStream.writeBytes(twoHyphens + boundary + lineEnd);
            outputStream.writeBytes("Content-Disposition: form-data; name=\"" + filefield + "\"; filename=\"" + q[idx] +"\"" + lineEnd);
            // outputStream.writeBytes("Content-Type: video/mp4" + lineEnd);
            outputStream.writeBytes("Content-Type: " + getMimeType(filepath) + lineEnd);

            outputStream.writeBytes("Content-Transfer-Encoding: binary" + lineEnd);
            outputStream.writeBytes(lineEnd);

            bytesAvailable = fileInputStream.available();
            bufferSize = Math.min(bytesAvailable, maxBufferSize);
            buffer = new byte[bufferSize];

            bytesRead = fileInputStream.read(buffer, 0, bufferSize);
            while(bytesRead > 0) {
                outputStream.write(buffer, 0, bufferSize);
                bytesAvailable = fileInputStream.available();
                bufferSize = Math.min(bytesAvailable, maxBufferSize);
                bytesRead = fileInputStream.read(buffer, 0, bufferSize);

            }
            outputStream.writeBytes(lineEnd);

            // Upload POST Data
//                String[] posts = post.split("&");
//                int max = posts.length;
//                for(int i=0; i<max;i++) {
//                        outputStream.writeBytes(twoHyphens + boundary + lineEnd);
//                        String[] kv = posts[i].split("=");
//                        outputStream.writeBytes("Content-Disposition: form-data; name=\"" + kv[0] + "\"" + lineEnd);
//                        outputStream.writeBytes("Content-Type: text/plain"+lineEnd);
//                        outputStream.writeBytes(lineEnd);
//                        outputStream.writeBytes(kv[1]);
//                        outputStream.writeBytes(lineEnd);
//                }
//
            Iterator iter = httpPara.entrySet().iterator();

            while (iter.hasNext()) {
                Map.Entry para = (Map.Entry) iter.next();
                outputStream.writeBytes(twoHyphens + boundary + lineEnd);
                outputStream.writeBytes("Content-Disposition: form-data; name=\"" + para.getKey().toString() + "\"" + lineEnd);
                outputStream.writeBytes("Content-Type: text/plain"+lineEnd);
                outputStream.writeBytes(lineEnd);
                outputStream.writeBytes(para.getValue().toString());
                outputStream.writeBytes(lineEnd);
            }
            outputStream.writeBytes(twoHyphens + boundary + twoHyphens + lineEnd);
            outputStream.flush();
            outputStream.close();

            ResponseCode = connection.getResponseCode();
            inputStream = connection.getInputStream();
            result = this.convertStreamToString(inputStream);

            fileInputStream.close();
            inputStream.close();


            return new ResponseObject(ResponseCode, result);
        } catch(Exception e) {
            Log.e("MultipartRequest", "Multipart Form Upload Error");

            e.printStackTrace();

            return null;
        }
    }
}
*/
