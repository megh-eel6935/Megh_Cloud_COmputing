package com.example.salman.login;

import android.app.Activity;
import android.content.ContentResolver;
import android.content.Context;
import android.content.Intent;
import android.content.res.Configuration;
import android.database.CharArrayBuffer;
import android.database.ContentObserver;
import android.database.Cursor;
import android.database.DataSetObserver;
import android.net.Uri;
import android.os.AsyncTask;
import android.os.Bundle;
import android.provider.MediaStore;
import android.provider.OpenableColumns;
import android.support.v4.widget.DrawerLayout;
import android.support.v7.app.ActionBarActivity;
import android.support.v7.app.ActionBarDrawerToggle;
import android.util.Log;
import android.view.LayoutInflater;
import android.view.MenuItem;
import android.view.View;
import android.view.ViewGroup;
import android.widget.AdapterView;
import android.widget.ArrayAdapter;
import android.widget.BaseAdapter;
import android.widget.ImageView;
import android.widget.ListView;
import android.widget.TextView;
import android.widget.Toast;

import com.example.salman.login.CustomObjects.GroupData;
import com.example.salman.login.CustomObjects.Groups;
import com.example.salman.login.CustomObjects.SelectedContent;
import com.example.salman.login.CustomObjects.UserMessage;

import java.util.List;


/**
 * Created by Ahmad Salman Saqib on 11/27/14.
 */
public class DrawerActivity extends ActionBarActivity {

    private static final int READ_REQUEST_CODE = 42;

    SelectedContent selectedContent;

    private DrawerLayout drawerLayout;
    private ListView drawerLeftListView;
    private ListView drawer_baseLayout_listView;
    private ListView drawerRightListView;

    private ActionBarDrawerToggle drawerLeftListener;

    private AdapterLeftDrawerList adapterLeftDrawerList;
    private AdapterMessage adapterMessage;

    private List<Groups> groupsList;
    private List<GroupData> groupsDataList;

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_drawer);

        drawerLayout = (DrawerLayout) findViewById(R.id.drawer_activityID);
        drawerLeftListView = (ListView) findViewById(R.id.drawer_left_slider_listView);
        drawer_baseLayout_listView = (ListView) findViewById(R.id.drawer_baseLayout_listView);
        drawerRightListView = (ListView) findViewById(R.id.drawer_right_slider_listView);

        adapterLeftDrawerList = new AdapterLeftDrawerList(this);
        drawerLeftListView.setAdapter(adapterLeftDrawerList);


        //adapterMessage = new AdapterMessage(this);

        drawerLeftListView.setOnItemClickListener(new LeftDrawerItemClickListener());
        drawerRightListView.setOnItemClickListener(new RightDrawerItemClickListener());

        drawerLeftListener = new ActionBarDrawerToggle(this, drawerLayout,
                R.string.drawer_open, R.string.drawer_close) {

            @Override
            public void onDrawerOpened(View drawerView) {
                //super.onDrawerOpened(drawerView);

            }

            @Override
            public void onDrawerClosed(View drawerView) {
                //super.onDrawerClosed(drawerView);

            }
        };

        drawerLayout.setDrawerListener(drawerLeftListener);
        getSupportActionBar().setHomeButtonEnabled(true);
        getSupportActionBar().setDisplayHomeAsUpEnabled(true);
    }

    @Override
    protected void onSaveInstanceState(Bundle outState) {
        super.onSaveInstanceState(outState);
        //TODO:outState.put("lastUsedAdapter",drawer_baseLayout_listView.getAdapter());
    }

    @Override
    protected void onRestoreInstanceState(Bundle savedInstanceState) {
        super.onRestoreInstanceState(savedInstanceState);
    }

    @Override
    protected void onResume() {
        super.onResume();
    }

    @Override
    public boolean onOptionsItemSelected(MenuItem item) {
        if (drawerLeftListener.onOptionsItemSelected(item)) {
            return true;
        }
        return super.onOptionsItemSelected(item);
    }

    @Override
    public void onConfigurationChanged(Configuration newConfig) {
        super.onConfigurationChanged(newConfig);
        drawerLeftListener.onConfigurationChanged(newConfig);
    }

    @Override
    protected void onPostCreate(Bundle savedInstanceState) {
        super.onPostCreate(savedInstanceState);
        drawerLeftListener.syncState();
    }

    @Override
    protected void onActivityResult(int requestCode, int resultCode,
                                    Intent resultData) {
        //Log.v("result_data",resultData.toString());
        if (requestCode == READ_REQUEST_CODE && resultCode == Activity.RESULT_OK) {

            Uri uri;
            if (resultData != null) {
                //Log.v("result_data 2", "not null");
                uri = resultData.getData();
                //Log.i("mara kha", "Uri: " + uri.toString());
                //showImage(uri);
                String metaData = dumpFileMetaData(uri);

                String realPath = getRealPathFromURI(uri);
                Log.v("getRealPath", realPath);
                Database.database.put("currentUri", realPath);
                Log.v("URIatSource", Database.database.get("currentUri"));
                Database.database.put("currentFileName", metaData);
            }
        }
        try {
            boolean success = UploadFile.handleFile(Database.database.get("currentUri"));
            if (success) {
                Toast.makeText(getApplicationContext(), "success", Toast.LENGTH_SHORT).show();
            }
            else {Toast.makeText(getApplicationContext(),"Failure", Toast.LENGTH_SHORT).show();
            }
        } catch (Exception e) {
            e.printStackTrace();
        }
    }



    /*
    * DrawerActivity custom methods
    * */

    public void setTitle(String title) {
        getSupportActionBar().setTitle(title);
    }



    /*
    * DrawerActivity click listeners
    * */

    private class RightDrawerItemClickListener implements AdapterView.OnItemClickListener {

        @Override
        public void onItemClick(AdapterView<?> parent, View view, int position, long id) {
            if (position == 0) {
                Intent intent = new Intent(Intent.ACTION_OPEN_DOCUMENT);
                intent.addCategory(Intent.CATEGORY_OPENABLE);
                intent.setType("*/*");
                startActivityForResult(intent, READ_REQUEST_CODE);
                Log.v("1","1");

                //Intent chooseFile = new Intent(getApplicationContext(), FileChooser.class);
                //startActivity(chooseFile);
                //BTGetSelectedItemHandler itemClicked = new BTGetSelectedItemHandler();
                //itemClicked.execute(position);
            }
            //Log.v("3","3");
            drawerRightListView.setItemChecked(position, true);
            drawerLayout.closeDrawer(drawerRightListView);
        }

        private class BTGetSelectedItemHandler extends AsyncTask<Integer, String, Boolean> {

            @Override
            protected Boolean doInBackground(Integer... params) {

                boolean success = false;
                if (params[0] == 0) {
                    try {

                        success = UploadFile.handleFile(Database.database.get("currentUri"));
                        return success;

                    } catch (Exception e) {
                        e.printStackTrace();
                    }

                }
                return success;
            }

            @Override
            protected void onPostExecute(Boolean b) {
                super.onPostExecute(b);
                if (b) {
                    Toast.makeText(getApplicationContext(),"success", Toast.LENGTH_SHORT).show();
                }
                else {Toast.makeText(getApplicationContext(),"Failure", Toast.LENGTH_SHORT).show();}
            }
        }

    }

    private class LeftDrawerItemClickListener implements AdapterView.OnItemClickListener {

        @Override
        public void onItemClick(AdapterView<?> parent, View view, int position, long id) {
            BTGetSelectedItemHandler itemClicked = new BTGetSelectedItemHandler();
            itemClicked.execute(position);

            drawerLeftListView.setItemChecked(position, true);
            setTitle(adapterLeftDrawerList.leftDrawerOptions[position]);
            drawerLayout.closeDrawer(drawerLeftListView);
        }

        private class BTGetSelectedItemHandler extends AsyncTask<Integer, String, SelectedContent> {

            @Override
            protected SelectedContent doInBackground(Integer... params) {
                try {
                    String content = UserFunctions.getSelectedListItemData(params[0]);
                    int position = params[0];

                    selectedContent = new SelectedContent();
                    selectedContent.setContent(content);
                    selectedContent.setPosition(position);
                    return selectedContent;

                } catch (Exception e) {
                    e.printStackTrace();
                    return null;
                }
            }

            @Override
            protected void onPostExecute(SelectedContent s) {
                if (s.getPosition() == 0) {
                    List<UserMessage> messageList = JsonCustomParser.readAndParseJSONMessages(s.getContent());
                    adapterMessage = new AdapterMessage(DrawerActivity.this,
                            R.layout.internallayout_drawer_base_item, messageList);

                    drawer_baseLayout_listView.setAdapter(adapterMessage);

                } else if (s.getPosition() == 1) {
                    groupsList = JsonCustomParser.readAndParseJSONGroups(s.getContent());
                    AdapterGroups adapterGroups = new AdapterGroups(DrawerActivity.this,
                            R.layout.internallayout_drawer_base_item, groupsList);

                    drawer_baseLayout_listView.setAdapter(adapterGroups);


                    drawer_baseLayout_listView.setOnItemClickListener(new GroupsItemClickListener());

                }
            }
        }
    }

    private class GroupsItemClickListener implements AdapterView.OnItemClickListener {

        @Override
        public void onItemClick(AdapterView<?> parent, View view, int position, long id) {

            BTGetSelectedGroupHandler groupClicked = new BTGetSelectedGroupHandler();
            groupClicked.execute(position);

            drawer_baseLayout_listView.setItemChecked(position, true);
            setTitle(groupsList.get(position).getGroup_name());
        }

        private class BTGetSelectedGroupHandler extends AsyncTask<Integer, String, String> {

            String grpID;
            String content;

            @Override
            protected String doInBackground(Integer... params) {
                grpID = groupsList.get(params[0]).getGroup_id();
                try {
                    content = UserFunctions.getSelectedGroupData(grpID);
                    return content;
                } catch (Exception e) {
                    e.printStackTrace();
                    return null;
                }
            }

            @Override
            protected void onPostExecute(String s) {
                groupsDataList = JsonCustomParser.readAndParseJSONGroupData(s);
                AdapterGroupData adapterGroupData = new AdapterGroupData(DrawerActivity.this,
                        R.layout.internallayout_drawer_base_item, groupsDataList);

                drawer_baseLayout_listView.setAdapter(adapterGroupData);
            }
        }
    }

    public String dumpFileMetaData(Uri uri) {

        // The query, since it only applies to a single document, will only return
        // one row. There's no need to filter, sort, or select fields, since we want
        // all fields for one document.
        Cursor cursor = new Cursor() {
            @Override
            public int getCount() {
                return 0;
            }

            @Override
            public int getPosition() {
                return 0;
            }

            @Override
            public boolean move(int offset) {
                return false;
            }

            @Override
            public boolean moveToPosition(int position) {
                return false;
            }

            @Override
            public boolean moveToFirst() {
                return false;
            }

            @Override
            public boolean moveToLast() {
                return false;
            }

            @Override
            public boolean moveToNext() {
                return false;
            }

            @Override
            public boolean moveToPrevious() {
                return false;
            }

            @Override
            public boolean isFirst() {
                return false;
            }

            @Override
            public boolean isLast() {
                return false;
            }

            @Override
            public boolean isBeforeFirst() {
                return false;
            }

            @Override
            public boolean isAfterLast() {
                return false;
            }

            @Override
            public int getColumnIndex(String columnName) {
                return 0;
            }

            @Override
            public int getColumnIndexOrThrow(String columnName) throws IllegalArgumentException {
                return 0;
            }

            @Override
            public String getColumnName(int columnIndex) {
                return null;
            }

            @Override
            public String[] getColumnNames() {
                return new String[0];
            }

            @Override
            public int getColumnCount() {
                return 0;
            }

            @Override
            public byte[] getBlob(int columnIndex) {
                return new byte[0];
            }

            @Override
            public String getString(int columnIndex) {
                return null;
            }

            @Override
            public void copyStringToBuffer(int columnIndex, CharArrayBuffer buffer) {

            }

            @Override
            public short getShort(int columnIndex) {
                return 0;
            }

            @Override
            public int getInt(int columnIndex) {
                return 0;
            }

            @Override
            public long getLong(int columnIndex) {
                return 0;
            }

            @Override
            public float getFloat(int columnIndex) {
                return 0;
            }

            @Override
            public double getDouble(int columnIndex) {
                return 0;
            }

            @Override
            public int getType(int columnIndex) {
                return 0;
            }

            @Override
            public boolean isNull(int columnIndex) {
                return false;
            }

            @Override
            public void deactivate() {

            }

            @Override
            public boolean requery() {
                return false;
            }

            @Override
            public void close() {

            }

            @Override
            public boolean isClosed() {
                return false;
            }

            @Override
            public void registerContentObserver(ContentObserver observer) {

            }

            @Override
            public void unregisterContentObserver(ContentObserver observer) {

            }

            @Override
            public void registerDataSetObserver(DataSetObserver observer) {

            }

            @Override
            public void unregisterDataSetObserver(DataSetObserver observer) {

            }

            @Override
            public void setNotificationUri(ContentResolver cr, Uri uri) {

            }

            @Override
            public Uri getNotificationUri() {
                return null;
            }

            @Override
            public boolean getWantsAllOnMoveCalls() {
                return false;
            }

            @Override
            public Bundle getExtras() {
                return null;
            }

            @Override
            public Bundle respond(Bundle extras) {
                return null;
            }
        };

        cursor = this.getContentResolver()
                .query(uri, null, null, null, null, null);

        try {
            // moveToFirst() returns false if the cursor has 0 rows.  Very handy for
            // "if there's anything to look at, look at it" conditionals.
            if (cursor != null && cursor.moveToFirst()) {

                // Note it's called "Display Name".  This is
                // provider-specific, and might not necessarily be the file name.
                String displayName = cursor.getString(
                        cursor.getColumnIndex(OpenableColumns.DISPLAY_NAME));
                //Log.i("mara kha", "Display Name: " + displayName);

                int sizeIndex = cursor.getColumnIndex(OpenableColumns.SIZE);
                // If the size is unknown, the value stored is null.  But since an
                // int can't be null in Java, the behavior is implementation-specific,
                // which is just a fancy term for "unpredictable".  So as
                // a rule, check if it's null before assigning to an int.  This will
                // happen often:  The storage API allows for remote files, whose
                // size might not be locally known.
                String size;
                if (!cursor.isNull(sizeIndex)) {
                    // Technically the column stores an int, but cursor.getString()
                    // will do the conversion automatically.
                    size = cursor.getString(sizeIndex);
                } else {
                    size = "Unknown";
                }
                //Log.i("mara kha", "Size: " + size);
                return  displayName;
            }
        } finally {
            cursor.close();
            return null;
        }
    }

    public String getRealPathFromURI(Uri contentUri) {
        Cursor cursor = null;
        try {
            Log.v("getRealPath", contentUri.toString() );
            String[] proj = {MediaStore.Images.Media.DATA};

            cursor = this.getContentResolver().query(contentUri, proj, null, null, null);

            int column_index = cursor.getColumnIndexOrThrow(MediaStore.Images.Media.DATA);

            cursor.moveToFirst();
            String rtrn = cursor.getString(column_index);
            Log.v("realPath",rtrn);
            return rtrn;
        } finally {
            if (cursor != null) {

                cursor.close();
            }
        }
    }
}





/*
* Adapters
* */

class AdapterLeftDrawerList extends BaseAdapter {

    private Context context;
    String[] leftDrawerOptions;
    int[] images = {R.drawable.ic_home, R.drawable.ic_people, R.drawable.ic_groups};


    public AdapterLeftDrawerList(Context context) {
        this.context = context;
        leftDrawerOptions = context.getResources().getStringArray(R.array.nav_drawer_left_items);
    }

    @Override
    public int getCount() {
        return leftDrawerOptions.length;
    }

    @Override
    public Object getItem(int position) {
        return leftDrawerOptions[position];
    }

    @Override
    public long getItemId(int position) {
        return position;
    }

    @Override
    public View getView(int position, View convertView, ViewGroup parent) {
        View row;
        if (convertView == null) {
            LayoutInflater inflater = (LayoutInflater) context.getSystemService(Context.LAYOUT_INFLATER_SERVICE);
            row = inflater.inflate(R.layout.internallayout_drawer_listitem, parent, false);
        } else {
            row = convertView;
        }
        TextView titleTextView = (TextView) row.findViewById(R.id.internalLayout_listItem_title_textView);
        ImageView titleImageView = (ImageView) row.findViewById(R.id.internalLayout_listItem_icon_imageView);

        titleTextView.setText(leftDrawerOptions[position]);
        titleImageView.setImageResource(images[position]);

        return row;
    }
}


class AdapterMessage extends ArrayAdapter<UserMessage> {

    private Context context;
    private List<UserMessage> messageList;

    public AdapterMessage(Context context, int resource, List<UserMessage> objects) {
        super(context, resource, objects);
        this.context = context;
        this.messageList = objects;
    }

    public View getView(int position, View convertView, ViewGroup parent) {
        LayoutInflater inflater =
                (LayoutInflater) context.getSystemService(Activity.LAYOUT_INFLATER_SERVICE);

        View view = inflater.inflate(R.layout.internallayout_drawer_base_item, parent, false);

        UserMessage msg = messageList.get(position);
        TextView tv = (TextView) view.findViewById(R.id.internalLayout_baseItem_title_textView);
        tv.setText(msg.getContent());


        ImageView img = (ImageView) view.findViewById(R.id.internalLayout_baseItem_icon_imageView);

        if (msg.getContent_type().equals("text")) {
            img.setImageResource(R.drawable.ic_message);
        } else if (msg.getContent_type().equals("file")) {
            img.setImageResource(R.drawable.ic_file);
        } else if (msg.getContent_type().equals("url")) {
            img.setImageResource(R.drawable.ic_url);
        }
        return view;
    }
}

class AdapterGroups extends ArrayAdapter<Groups> {

    private Context context;
    private List<Groups> groupsList;

    public AdapterGroups(Context context, int resource, List<Groups> objects) {
        super(context, resource, objects);
        this.context = context;
        this.groupsList = objects;
    }

    public View getView(int position, View convertView, ViewGroup parent) {
        LayoutInflater inflater =
                (LayoutInflater) context.getSystemService(Activity.LAYOUT_INFLATER_SERVICE);

        View view = inflater.inflate(R.layout.internallayout_drawer_base_item, parent, false);

        Groups grp = groupsList.get(position);
        TextView tv = (TextView) view.findViewById(R.id.internalLayout_baseItem_title_textView);
        tv.setText(grp.getGroup_name());


        ImageView img = (ImageView) view.findViewById(R.id.internalLayout_baseItem_icon_imageView);

        img.setImageResource(R.drawable.ic_groups);

        return view;

    }
}

class AdapterGroupData extends ArrayAdapter<GroupData> {

    private Context context;
    private List<GroupData> groupDataList;

    public AdapterGroupData(Context context, int resource, List<GroupData> objects) {
        super(context, resource, objects);
        this.context = context;
        this.groupDataList = objects;
    }

    public View getView(int position, View convertView, ViewGroup parent) {
        LayoutInflater inflater =
                (LayoutInflater) context.getSystemService(Activity.LAYOUT_INFLATER_SERVICE);

        View view = inflater.inflate(R.layout.internallayout_drawer_base_item, parent, false);

        GroupData grpData = groupDataList.get(position);
        TextView tv = (TextView) view.findViewById(R.id.internalLayout_baseItem_title_textView);
        tv.setText(grpData.getContent());


        ImageView img = (ImageView) view.findViewById(R.id.internalLayout_baseItem_icon_imageView);

        if (grpData.getContent_type().equals("text")) {
            img.setImageResource(R.drawable.ic_message);
        } else if (grpData.getContent_type().equals("file")) {
            img.setImageResource(R.drawable.ic_file);
        } else if (grpData.getContent_type().equals("url")) {
            img.setImageResource(R.drawable.ic_url);
        }
        return view;
    }
}


