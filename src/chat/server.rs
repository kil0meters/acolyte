use std::collections::HashMap;

use actix::prelude::{Context, Handler, Recipient};
use actix::Actor;
use rand::{self, rngs::ThreadRng, Rng};
use serde_json;

use super::message_types::*;

pub struct Client {
    username: String,
    authentication_level: String,
}

pub struct Server {
    sessions: HashMap<usize, Recipient<Broadcast>>,
    rng: ThreadRng,
}

impl Default for Server {
    fn default() -> Server {
        Server {
            sessions: HashMap::new(),
            rng: rand::thread_rng(),
        }
    }
}

impl Server {
    fn broadcast_message(&self, message: String) {
        for addr in self.sessions.values() {
            let _ = addr.do_send(Broadcast(message.to_owned()));
        }
    }
}

impl Actor for Server {
    type Context = Context<Self>;
}

impl Handler<ChatMessage> for Server {
    type Result = ();

    fn handle(&mut self, msg: ChatMessage, _: &mut Context<Self>) {
        let m = serde_json::to_string(&msg).unwrap();

        self.broadcast_message(m);
    }
}

impl Handler<Connect> for Server {
    type Result = usize;

    fn handle(&mut self, msg: Connect, _: &mut Context<Self>) -> Self::Result {
        let id = self.rng.gen::<usize>();
        self.sessions.insert(id, msg.addr);

        id
    }
}

impl Handler<Disconnect> for Server {
    type Result = ();

    fn handle(&mut self, msg: Disconnect, _: &mut Context<Self>) {
        self.sessions.remove(&msg.id);
    }
}
