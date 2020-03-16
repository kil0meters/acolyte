use std::time;

use actix::*;
use actix_web_actors::ws;
use uuid::Uuid;

use super::message_types::*;
use super::server::Server;

use crate::auth::permissions::AuthLevel;

const HEARTBEAT_INTERVAL: time::Duration = time::Duration::from_secs(5);
const CLIENT_TIMEOUT: time::Duration = time::Duration::from_secs(10);

pub struct Client {
    pub id: usize,

    pub username: String,
    pub auth_level: AuthLevel,
    // If client doesn't send a ping every 10 seconds
    // it gets absolutely murdered
    pub hb: time::Instant,

    // Address of server connection
    pub conn: Addr<Server>,
}

impl Actor for Client {
    type Context = ws::WebsocketContext<Self>;

    fn started(&mut self, ctx: &mut Self::Context) {
        // Pings server
        self.start_heartbeat(ctx);

        let addr = ctx.address();
        self.conn
            .send(Connect {
                addr: addr.recipient(),
            })
            .into_actor(self)
            .then(|res, act, ctx| {
                match res {
                    Ok(res) => act.id = res,
                    // if the server isn't working just stop the connection
                    _ => ctx.stop(),
                }
                actix::fut::ready(())
            })
            .wait(ctx);
    }
}

impl Handler<Broadcast> for Client {
    type Result = ();

    // simply send broadcasts to the client
    fn handle(&mut self, msg: Broadcast, ctx: &mut Self::Context) {
        ctx.text(msg.0);
    }
}

impl StreamHandler<Result<ws::Message, ws::ProtocolError>> for Client {
    fn handle(&mut self, msg: Result<ws::Message, ws::ProtocolError>, ctx: &mut Self::Context) {
        let msg = match msg {
            Err(_) => {
                ctx.stop();
                return;
            }
            Ok(msg) => msg,
        };

        match msg {
            ws::Message::Ping(msg) => {
                self.hb = time::Instant::now();
                ctx.pong(&msg);
            }
            ws::Message::Pong(_) => {
                self.hb = time::Instant::now();
            }
            ws::Message::Text(text) => {
                let m = text.trim();
                println!("{}", m);
                if m.starts_with('/') {
                } else if &self.username != "ANON" {
                    self.conn.do_send(ChatMessage {
                        username: self.username.clone(),
                        id: Uuid::new_v4(),
                        date: time::SystemTime::now(),
                        text,
                    });
                }
            }
            // ws::Message::Binary(bin) => ctx.binary(bin),
            // ws::Message::Close(msg) => {}
            _ => (),
        }
    }
}

impl Client {
    fn start_heartbeat(&self, ctx: &mut ws::WebsocketContext<Self>) {
        ctx.run_interval(HEARTBEAT_INTERVAL, |act, ctx| {
            // If the timer runs out, drpo the connection from the pool
            if time::Instant::now().duration_since(act.hb) > CLIENT_TIMEOUT {
                act.conn.do_send(Disconnect { id: act.id });

                ctx.stop();

                return;
            }

            ctx.ping(b"");
        });
    }
}
