use rand::Rng;

pub fn hash_password(password: String) -> Result<String, argon2::Error> {
    let salt: [u8; 32] = rand::thread_rng().gen();
    let config = argon2::Config::default();

    argon2::hash_encoded(password.as_bytes(), &salt, &config)
}

pub fn verify_password(password: String, password_hash: String) -> Result<bool, argon2::Error> {
    argon2::verify_encoded(&password_hash, password.as_bytes())
}
