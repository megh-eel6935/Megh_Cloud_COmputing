package com.example.salman.login;

import android.content.Context;
import android.content.res.Configuration;
import android.os.AsyncTask;
import android.os.Bundle;
import android.support.v4.widget.DrawerLayout;
import android.support.v7.app.ActionBarActivity;
import android.support.v7.app.ActionBarDrawerToggle;
import android.view.LayoutInflater;
import android.view.MenuItem;
import android.view.View;
import android.view.ViewGroup;
import android.widget.AdapterView;
import android.widget.BaseAdapter;
import android.widget.ImageView;
import android.widget.ListView;
import android.widget.TextView;
import android.widget.Toast;

/**
 * Created by salman on 11/27/14.
 */
public class DrawerActivity extends ActionBarActivity {

    private DrawerLayout drawerLayout;
    private ListView drawerLeftListView;
    private ListView drawerRightListView;

    private ActionBarDrawerToggle drawerLeftListener;
    private MyAdapter myAdapter;

    private TextView testTextView;

    @Override
    protected  void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_drawer);

        drawerLayout = (DrawerLayout) findViewById(R.id.drawer_activityID);
        drawerLeftListView = (ListView) findViewById(R.id.drawer_left_slider_listView);
        drawerRightListView = (ListView) findViewById(R.id.drawer_right_slider_listView);

        //////////////////
        testTextView = (TextView) findViewById(R.id.drawer_insideTextView);
        /////////////////

        myAdapter = new MyAdapter(this);
        drawerLeftListView.setAdapter(myAdapter);

        drawerLeftListView.setOnItemClickListener(new ItemClickListener());

        drawerLeftListener = new ActionBarDrawerToggle(this, drawerLayout,
                R.string.drawer_open, R.string.drawer_close){

            @Override
            public void onDrawerOpened(View drawerView) {
                //super.onDrawerOpened(drawerView);
                Toast.makeText(DrawerActivity.this, "Drawer Opened", Toast.LENGTH_SHORT).show();
            }

            @Override
            public void onDrawerClosed(View drawerView) {
                //super.onDrawerClosed(drawerView);
                Toast.makeText(DrawerActivity.this, "Drawer Closed", Toast.LENGTH_SHORT).show();
            }
        };

        drawerLayout.setDrawerListener(drawerLeftListener);
        getSupportActionBar().setHomeButtonEnabled(true);
        getSupportActionBar().setDisplayHomeAsUpEnabled(true);
    }

    @Override
    public boolean onOptionsItemSelected(MenuItem item) {
        if (drawerLeftListener.onOptionsItemSelected(item)){
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

    private class ItemClickListener implements AdapterView.OnItemClickListener{

        @Override
        public void onItemClick(AdapterView<?> parent, View view, int position, long id) {

            String content = null;

            Toast.makeText(DrawerActivity.this, myAdapter.leftDrawerOptions[position] + " was selected",
                    Toast.LENGTH_SHORT).show();

            BTGetSelectedItemFunctionalityHandler itemClicked = new BTGetSelectedItemFunctionalityHandler();
            itemClicked.execute(position);

            selectMethod(position);
            drawerLayout.closeDrawer(drawerLeftListView);
        }

        private class BTGetSelectedItemFunctionalityHandler extends AsyncTask<Integer, String, String> {

            @Override
            protected String doInBackground(Integer... params) {
                try {
                    String content = UserFunctions.getSelectedListItemFunctionality(params[0]);
                    return content;

                } catch (Exception e) {
                    e.printStackTrace();
                    return null;
                }
            }

            @Override
            protected void onPostExecute(String s) {
                testTextView.append(s);
            }
        }
    }

    public void selectMethod(int position) {
        drawerLeftListView.setItemChecked(position, true);
        setTitle(myAdapter.leftDrawerOptions[position]);
    }

    public void setTitle(String title) {
        getSupportActionBar().setTitle(title);
    }
}

class MyAdapter extends BaseAdapter {

    private Context context;
    String[] leftDrawerOptions;
    int[] images = {R.drawable.ic_home, R.drawable.ic_people,R.drawable.ic_groups};


    public MyAdapter(Context context){
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
        if (convertView == null){
            LayoutInflater inflater = (LayoutInflater) context.getSystemService(Context.LAYOUT_INFLATER_SERVICE);
            row = inflater.inflate(R.layout.internallayout_drawer_listitem, parent, false);
        }
        else{
            row = convertView;
        }
        TextView titleTextView = (TextView) row.findViewById(R.id.internalLayout_listItem_title_textView);
        ImageView titleImageView = (ImageView) row.findViewById(R.id.internalLayout_listItem_icon_imageView);

        titleTextView.setText(leftDrawerOptions[position]);
        titleImageView.setImageResource(images[position]);

        return row;
    }
}