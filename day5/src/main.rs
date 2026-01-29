use actix_web::{web, App, HttpServer, HttpResponse, Responder};
use deadpool_postgres::{Config, ManagerConfig, RecyclingMethod, Runtime, Pool};
use tokio_postgres::NoTls;
use serde::Serialize;
use std::env;

#[derive(Serialize)]
struct Status {
    status: String,
}

#[derive(Serialize)]
struct DbInfo {
    version: String,
}

async fn index() -> impl Responder {
    HttpResponse::Ok().json(Status { status: "Rust API is running! 🦀".to_string() })
}

async fn db_check(pool: web::Data<Pool>) -> impl Responder {
    let client = match pool.get().await {
        Ok(client) => client,
        Err(e) => return HttpResponse::InternalServerError().json(Status { status: format!("Failed to get connection: {}", e) }),
    };

    let rows = match client.query("SELECT version()", &[]).await {
        Ok(rows) => rows,
        Err(e) => return HttpResponse::InternalServerError().json(Status { status: format!("Query failed: {}", e) }),
    };

    if rows.is_empty() {
        return HttpResponse::InternalServerError().json(Status { status: "No rows returned".to_string() });
    }

    let value: &str = rows[0].get(0);
    HttpResponse::Ok().json(DbInfo { version: value.to_string() })
}

#[actix_web::main]
async fn main() -> std::io::Result<()> {
    let mut cfg = Config::new();
    cfg.host = Some(env::var("PG_HOST").unwrap_or("postgres".to_string()));
    cfg.user = Some(env::var("PG_USER").unwrap_or("postgres".to_string()));
    cfg.password = Some(env::var("PG_PASSWORD").unwrap_or("password".to_string()));
    cfg.dbname = Some(env::var("PG_DBNAME").unwrap_or("appdb".to_string()));
    
    // Safety fallback for port
    if let Ok(port) = env::var("PG_PORT") {
        if let Ok(p) = port.parse::<u16>() {
            cfg.port = Some(p);
        }
    }

    cfg.manager = Some(ManagerConfig { recycling_method: RecyclingMethod::Fast });

    let pool = cfg.create_pool(Some(Runtime::Tokio1), NoTls).expect("Failed to create pool");

    println!("Server starting at http://0.0.0.0:8080");

    HttpServer::new(move || {
        App::new()
            .app_data(web::Data::new(pool.clone()))
            .route("/", web::get().to(index))
            .route("/db", web::get().to(db_check))
    })
    .bind(("0.0.0.0", 8080))?
    .run()
    .await
}
