-- Tokens are stateless, so logging out or changing a password cannot invalidate
-- one on its own. This column records the moment every token issued earlier stops
-- being accepted; the auth middleware compares it against the token's iat claim.
-- Safe to run on a populated table: the column is nullable, and a NULL means
-- "no invalidation recorded yet".
ALTER TABLE "public"."users" ADD COLUMN IF NOT EXISTS "tokens_valid_after" timestamp;

COMMENT ON COLUMN "public"."users"."tokens_valid_after" IS 'tokens issued before this moment are rejected; set on logout and password change.';
