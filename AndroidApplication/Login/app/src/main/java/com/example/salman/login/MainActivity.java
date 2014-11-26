package com.example.salman.login;

import android.content.Context;
import android.content.SharedPreferences;
import android.content.pm.PackageInfo;
import android.content.pm.PackageManager;
import android.net.ConnectivityManager;
import android.net.NetworkInfo;
import android.os.AsyncTask;
import android.os.Bundle;
import android.support.v7.app.ActionBarActivity;
import android.text.method.ScrollingMovementMethod;
import android.util.Log;
import android.view.Menu;
import android.view.MenuItem;
import android.view.View;
import android.view.View.OnClickListener;
import android.widget.Button;
import android.widget.EditText;
import android.widget.Scroller;
import android.widget.TextView;
import android.widget.Toast;

import com.google.android.gms.common.ConnectionResult;
import com.google.android.gms.common.GooglePlayServicesUtil;
import com.google.android.gms.gcm.GoogleCloudMessaging;

import java.io.IOException;
import java.util.concurrent.atomic.AtomicInteger;

public class MainActivity extends ActionBarActivity {

    Button login_button;
    Button register_button;

    TextView temp_output;
    EditText login_email;
    EditText login_password;



    /////////////////////////////////////////////////////////////////////////////////////

    public static final String EXTRA_MESSAGE = "message";
    public static final String PROPERTY_REG_ID = "registration_id";
    private static final String PROPERTY_APP_VERSION = "appVersion";
    private static final int PLAY_SERVICES_RESOLUTION_REQUEST = 9000;


    String SENDER_ID = "492901091946";
    TextView mDisplay;
    GoogleCloudMessaging gcm;
    AtomicInteger msgId = new AtomicInteger();
    SharedPreferences prefs;
    Context context;

    String regid;


    static final String TAG = "GCM Megh";



    /////////////////////////////////////////////////////////////////////////////////////

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_main);

        login_button = (Button) findViewById(R.id.login_login_button);
        register_button = (Button) findViewById(R.id.login_register_button);

        login_email = (EditText) findViewById(R.id.login_email_editText);
        login_password = (EditText) findViewById(R.id.login_password_editText);

        temp_output = (TextView) findViewById(R.id.login_testOutput_textView);
        ///////////////////////////////////////////////////////////////////////////////////////||
        temp_output.setScroller(new Scroller(this)); //                                        ||
        temp_output.setVerticalScrollBarEnabled(true); //                                      ||
        temp_output.setMovementMethod(new ScrollingMovementMethod());//                        ||
        ///////////////////////////////////////////////////////////////////////////////////////||

        context = getApplicationContext();

        if (checkPlayServices()) {
            gcm = GoogleCloudMessaging.getInstance(this);
            regid = getRegistrationId(context);

            Log.v("RegistrationIDLog1","AAAAA:__"+regid + "__:AAAAA");

            if (regid.isEmpty()) {
                Log.i(TAG, "register in background commented out");

                RegisterInBackground registerInBackground = new RegisterInBackground();
                registerInBackground.execute();
            }
        } else {
            Log.i(TAG, "No valid Google Play Services APK found.");
        }


        login_button.setOnClickListener(new Login());
    }

    private class Login implements OnClickListener {

        @Override
        public void onClick(View v) {
            if (isOnline()) {
                requestLogin("http://104.131.126.89/mobilelogin");
            } else {
                Toast.makeText(v.getContext(),
                        "Network Not Available.\nCheck Network Connection",
                        Toast.LENGTH_LONG).show();
            }
        }

        private void requestLogin(String uri) {
            LoginRequestHandler LoginButton_clicked = new LoginRequestHandler();
            RequestPackage requestPackage = new RequestPackage();
            requestPackage.setUri(uri);
            requestPackage.setParam("email", login_email.getText().toString());
            requestPackage.setParam("password", login_password.getText().toString());
            Log.v("RegistrationIDLog2",regid);
            requestPackage.setParam("registration_id", regid);
            LoginButton_clicked.execute(requestPackage);
        }
    }

    @Override
    public boolean onCreateOptionsMenu(Menu menu) {
        // Inflate the menu; this adds items to the action bar if it is present.
        getMenuInflater().inflate(R.menu.menu_main, menu);
        return true;
    }

    @Override
    public boolean onOptionsItemSelected(MenuItem item) {
        // Handle action bar item clicks here. The action bar will
        // automatically handle clicks on the Home/Up button, so long
        // as you specify a parent activity in AndroidManifest.xml.
        int id = item.getItemId();

        //noinspection SimplifiableIfStatement
        if (id == R.id.action_settings) {
            return true;
        }

        return super.onOptionsItemSelected(item);
    }


    private class LoginRequestHandler extends AsyncTask<RequestPackage, String, String> {

        private final String LOG_TAG = MainActivity.class.getSimpleName();

        @Override
        protected void onPreExecute() {
            temp_output.append("PreExec" + "\n");
        }

        @Override
        protected String doInBackground(RequestPackage... params) {
            //String temp = params[0]+" : "+params[1];
            //return temp;
            String loginCookie;
            String test; /*TODO: remove this shit*/

            try {
                loginCookie = UserFunctions.sendAuthenticationPostMsg(params[0]);
                //return loginCookie;
                test = UserFunctions.sendGetMsg("http://104.131.126.89/getmessages", loginCookie);
                return test;
            } catch (Exception e) {
                e.printStackTrace();
                return null;
            }
        }

        @Override
        protected void onPostExecute(String s) {
            if (s == null) {
                Toast.makeText(MainActivity.this, "Cannot connect to web server", Toast.LENGTH_LONG).show();
                return;
            }
            else {
                temp_output.append(s+"\n");
            }
        }
    }


    protected boolean isOnline() {
        ConnectivityManager cm = (ConnectivityManager) getSystemService(Context.CONNECTIVITY_SERVICE);
        NetworkInfo netInfo = cm.getActiveNetworkInfo();
        if (netInfo != null && netInfo.isConnectedOrConnecting()) {
            return true;
        } else {
            return false;
        }
    }

    private boolean checkPlayServices() {
        int resultCode = GooglePlayServicesUtil.isGooglePlayServicesAvailable(this);
        if (resultCode != ConnectionResult.SUCCESS) {
            if (GooglePlayServicesUtil.isUserRecoverableError(resultCode)) {
                GooglePlayServicesUtil.getErrorDialog(resultCode, this,
                        PLAY_SERVICES_RESOLUTION_REQUEST).show();
            } else {
                Log.i(TAG, "This device is not supported.");
                finish();
            }
            return false;
        }
        return true;
    }

    //TODO: remove comments
    /**
     * Gets the current registration ID for application on GCM service.
     * <p>
     * If result is empty, the app needs to register.
     *
     * @return registration ID, or empty string if there is no existing
     *         registration ID.
     */
    private String getRegistrationId(Context context) {
        final SharedPreferences prefs = getGCMPreferences(context);
        String registrationId = prefs.getString(PROPERTY_REG_ID, "");
        if (registrationId.isEmpty()) {
            Log.i(TAG, "Registration not found.");
            return "";
        }
        // Check if app was updated; if so, it must clear the registration ID
        // since the existing regID is not guaranteed to work with the new
        // app version.
        int registeredVersion = prefs.getInt(PROPERTY_APP_VERSION, Integer.MIN_VALUE);
        int currentVersion = getAppVersion(context);
        if (registeredVersion != currentVersion) {
            Log.i(TAG, "App version changed.");
            return "";
        }
        return registrationId;
    }

    /**
     * @return Application's {@code SharedPreferences}.
     */
    private SharedPreferences getGCMPreferences(Context context) {
        // This sample app persists the registration ID in shared preferences, but
        // how you store the regID in your app is up to you.
        return getSharedPreferences(MainActivity.class.getSimpleName(),
                Context.MODE_PRIVATE);
    }

    /**
     * @return Application's version code from the {@code PackageManager}.
     */
    private static int getAppVersion(Context context) {
        try {
            PackageInfo packageInfo = context.getPackageManager().getPackageInfo(context.getPackageName(), 0);
            return packageInfo.versionCode;
        } catch (PackageManager.NameNotFoundException e) {
            // should never happen
            throw new RuntimeException("Could not get package name: " + e);
        }
    }

    /**
     * Registers the application with GCM servers asynchronously.
     * <p>
     * Stores the registration ID and app versionCode in the application's
     * shared preferences.
     */
    //TODO: registerInBackground implementation


    private class RegisterInBackground extends AsyncTask<Void,Void,String> {

        @Override
        protected String doInBackground(Void... params) {
            String msg = "";
            try {
                if (gcm == null) {
                    gcm = GoogleCloudMessaging.getInstance(context);
                }
                regid = gcm.register(SENDER_ID);
                msg = "Device registered, registration ID=" + regid;

                // You should send the registration ID to your server over HTTP,
                // so it can use GCM/HTTP or CCS to send messages to your app.
                // The request to your server should be authenticated if your app
                // is using accounts.

                //TODO:take care of the next line (sendRegistrationIdToBackend())
                //sendRegistrationIdToBackend();

                // For this demo: we don't need to send it because the device
                // will send upstream messages to a server that echo back the
                // message using the 'from' address in the message.

                // Persist the regID - no need to register again.
                storeRegistrationId(context, regid);
            } catch (IOException ex) {
                msg = "Error :" + ex.getMessage();
                // If there is an error, don't just keep trying to register.
                // Require the user to click a button again, or perform
                // exponential back-off.
            }
            return msg;
        }

        @Override
        protected void onPostExecute(String msg) {
            mDisplay.append(msg + "\n");
        }
    }

    /**
     * Stores the registration ID and app versionCode in the application's
     * {@code SharedPreferences}.
     *
     * @param context application's context.
     * @param regId registration ID
     */
    private void storeRegistrationId(Context context, String regId) {
        final SharedPreferences prefs = getGCMPreferences(context);
        int appVersion = getAppVersion(context);
        Log.i(TAG, "Saving regId on app version " + appVersion);
        SharedPreferences.Editor editor = prefs.edit();
        editor.putString(PROPERTY_REG_ID, regId);
        editor.putInt(PROPERTY_APP_VERSION, appVersion);
        editor.commit();
    }
}
