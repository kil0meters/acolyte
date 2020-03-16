use anyhow::Result;
use diesel::pg::PgConnection;
use diesel::prelude::*;

use crate::models::{Thread, User};
use crate::schema::threads;

const PAGE_LENGTH: i64 = 25;

pub fn get_hot_threads(page_number: i64, conn: &PgConnection) -> Result<Vec<Thread>> {
    let recent_threads = threads::table
        .select(threads::all_columns)
        .order_by(threads::created_at)
        .offset(page_number * PAGE_LENGTH)
        .limit(PAGE_LENGTH)
        .load::<Thread>(conn)?;

    debug!("threads: {:?}", recent_threads);

    Ok(recent_threads)
}

pub fn create_new_thread(
    author: User,
    title: String,
    body: String,
    link: String,
    conn: &PgConnection,
) -> Result<Thread> {
    let thread = Thread::new(author, title, body, link);

    thread.validate()?;
    Ok(diesel::insert_into(threads::table)
        .values(&thread)
        .get_result::<Thread>(conn)?)
}
