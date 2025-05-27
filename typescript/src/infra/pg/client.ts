import { Client } from "postgres";

export interface DatabaseConfig {
  host: string;
  port: number;
  database: string;
  username: string;
  password: string;
  ssl?: boolean;
}

export class DatabaseClient {
  private client: Client;
  private isConnected: boolean = false;

  constructor(private config: DatabaseConfig) {
    this.client = new Client({
      hostname: config.host,
      port: config.port,
      database: config.database,
      user: config.username,
      password: config.password,
      tls: config.ssl ? { enforce: true } : { enforce: false },
    });
  }

  async connect(): Promise<void> {
    if (!this.isConnected) {
      await this.client.connect();
      this.isConnected = true;
    }
  }

  async disconnect(): Promise<void> {
    if (this.isConnected) {
      await this.client.end();
      this.isConnected = false;
    }
  }

  async query<T>(text: string, params?: unknown[]): Promise<T[]> {
    await this.connect();
    const result = await this.client.queryObject<T>(text, params);
    return result.rows;
  }

  async queryRow<T>(text: string, params?: unknown[]): Promise<T | null> {
    const results = await this.query<T>(text, params);
    return results.length > 0 ? results[0] : null;
  }

  getClient(): Client {
    return this.client;
  }
}
