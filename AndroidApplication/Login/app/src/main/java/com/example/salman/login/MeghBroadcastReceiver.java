
package com.example.salman.login;

import android.content.Context;
import android.content.Intent;
import android.support.v4.content.WakefulBroadcastReceiver;

/**
 * Created by Ahmad Salman Saqib on 11/25/14.
 */


public class MeghBroadcastReceiver extends WakefulBroadcastReceiver {

    @Override
    public void onReceive(Context context, Intent intent) {
        // Explicitly specify that GcmIntentService will handle the intent.
        //ComponentName comp = new ComponentName(context.getPackageName(),
                //MeghIntentService.class.getName());
        // Start the service, keeping the device awake while it is launching.
        //startWakefulService(context, (intent.setComponent(comp)));
        //setResultCode(Activity.RESULT_OK);
    }
}

