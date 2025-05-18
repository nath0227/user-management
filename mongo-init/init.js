// Step 1: Create application user in the 'admin' database
db = db.getSiblingDB('admin');

db.createUser({
  user: "appusr",
  pwd: "Passw0rd",
  roles: [
    {
      role: "readWrite",
      db: "myapp"
    }
  ]
});

// Step 2: Switch to 'myapp' database for data setup
db = db.getSiblingDB('myapp');

// Ensure the 'users' collection exists
db.createCollection('users');

// Create a unique index on the email field
db.users.createIndex(
  { email: 1 },
  { unique: true }
);

// Insert an admin user
db.users.insertOne({
  name: "Admin",
  email: "admin@example.com",
  password: "$2a$10$jIEILNA1i4e57cjjfopkvOks3z22zZOVMaKvxmZU5V7C9I.9qL3FO", // bcrypt hash passwordstring
  createdAt: new Date()
});
