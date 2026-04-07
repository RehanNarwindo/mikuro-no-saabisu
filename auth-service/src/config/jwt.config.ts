const JwtConfig = () => ({
  jwt: {
    secret: process.env.JWT_SECRET,
    expiresIn: process.env.JWT_EXPIRES_IN || '10m',
  },
});

export default JwtConfig;
