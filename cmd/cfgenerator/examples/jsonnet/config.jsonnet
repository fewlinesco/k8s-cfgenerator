{
  database: {
    password: std.extVar('DATABASE_PASSWORD'),
    username: std.extVar('DATABASE_USERNAME'),
  },
  api: {
    address: '0.0.0.0:' + std.extVar('API_PORT'),
  },
}
