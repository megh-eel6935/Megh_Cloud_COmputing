package com.example.salman.login;

import android.content.Context;
import android.net.ConnectivityManager;
import android.net.NetworkInfo;
import android.os.AsyncTask;
import android.os.Bundle;
import android.support.v7.app.ActionBarActivity;
import android.text.method.ScrollingMovementMethod;
import android.view.Menu;
import android.view.MenuItem;
import android.view.View;
import android.view.View.OnClickListener;
import android.widget.Button;
import android.widget.EditText;
import android.widget.Scroller;
import android.widget.TextView;
import android.widget.Toast;

public class MainActivity extends ActionBarActivity {

    Button login_button;
    Button register_button;

    TextView temp_output;
    EditText login_email;
    EditText login_password;

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
        login_button.setOnClickListener(new Login());
    }

    private class Login implements OnClickListener {

        @Override
        public void onClick(View v) {
            if (isOnline()) {
                requestLogin("http://104.131.126.89/login");
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
            LoginButton_clicked.execute(requestPackage);
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
                loginCookie = UserFunctions.sendAuthenticationPost(params[0]);
                //return loginCookie;
                test = UserFunctions.sendGet("http://104.131.126.89/getmessages", loginCookie);
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
}
