BEGIN;

DO $$
BEGIN
    IF to_regclass('public.gomoku_games') IS NOT NULL THEN
        EXECUTE $sql$
            INSERT INTO games (id, game_type, player1_id, player2_id, winner, finished_at, created_at, updated_at, game_state)
            SELECT
                gg.id,
                'gomoku',
                gg.player1_id,
                gg.player2_id,
                gg.winner,
                gg.finished_at,
                gg.created_at,
                gg.updated_at,
                CASE
                    WHEN gg.game_state ? 'messages' THEN gg.game_state - 'messages'
                    ELSE gg.game_state
                END
            FROM gomoku_games gg
            ON CONFLICT (id) DO NOTHING
        $sql$;
    END IF;
END $$;

DO $$
BEGIN
    IF to_regclass('public.connectfour_games') IS NOT NULL THEN
        EXECUTE $sql$
            INSERT INTO games (id, game_type, player1_id, player2_id, winner, finished_at, created_at, updated_at, game_state)
            SELECT
                cg.id,
                'connectfour',
                cg.player1_id,
                cg.player2_id,
                cg.winner,
                cg.finished_at,
                cg.created_at,
                cg.updated_at,
                CASE
                    WHEN cg.game_state ? 'messages' THEN cg.game_state - 'messages'
                    ELSE cg.game_state
                END
            FROM connectfour_games cg
            ON CONFLICT (id) DO NOTHING
        $sql$;
    END IF;
END $$;

DO $$
BEGIN
    IF to_regclass('public.gomoku_messages') IS NOT NULL THEN
        EXECUTE $sql$
            INSERT INTO messages (game_id, game_type, sender_id, sender_name, content, sent_at)
            SELECT
                gm.game_id,
                'gomoku',
                gm.sender_id,
                gm.sender_name,
                gm.content,
                gm.sent_at
            FROM gomoku_messages gm
            WHERE EXISTS (
                SELECT 1
                FROM games g
                WHERE g.id = gm.game_id
                  AND g.game_type = 'gomoku'
            )
              AND NOT EXISTS (
                SELECT 1
                FROM messages m
                WHERE m.game_id = gm.game_id
                  AND m.game_type = 'gomoku'
                  AND m.sender_id = gm.sender_id
                  AND m.sender_name = gm.sender_name
                  AND m.content = gm.content
                  AND m.sent_at = gm.sent_at
            )
        $sql$;
    END IF;
END $$;

DO $$
BEGIN
    IF to_regclass('public.gomoku_games') IS NOT NULL THEN
        EXECUTE $sql$
            INSERT INTO messages (game_id, game_type, sender_id, sender_name, content, sent_at)
            SELECT
                gg.id,
                'gomoku',
                msg.value ->> 'senderID',
                msg.value ->> 'senderName',
                msg.value ->> 'content',
                COALESCE(NULLIF(msg.value ->> 'sentAt', '')::timestamptz, gg.finished_at)
            FROM gomoku_games gg
            CROSS JOIN LATERAL jsonb_array_elements(
                CASE
                    WHEN jsonb_typeof(gg.game_state -> 'messages') = 'array' THEN gg.game_state -> 'messages'
                    ELSE '[]'::jsonb
                END
            ) AS msg(value)
            WHERE NOT EXISTS (
                SELECT 1
                FROM messages m
                WHERE m.game_id = gg.id
                  AND m.game_type = 'gomoku'
                  AND m.sender_id = msg.value ->> 'senderID'
                  AND m.sender_name = msg.value ->> 'senderName'
                  AND m.content = msg.value ->> 'content'
                  AND m.sent_at = COALESCE(NULLIF(msg.value ->> 'sentAt', '')::timestamptz, gg.finished_at)
            )
        $sql$;
    END IF;
END $$;

DO $$
BEGIN
    IF to_regclass('public.connectfour_games') IS NOT NULL THEN
        EXECUTE $sql$
            INSERT INTO messages (game_id, game_type, sender_id, sender_name, content, sent_at)
            SELECT
                cg.id,
                'connectfour',
                msg.value ->> 'senderID',
                msg.value ->> 'senderName',
                msg.value ->> 'content',
                COALESCE(NULLIF(msg.value ->> 'sentAt', '')::timestamptz, cg.finished_at)
            FROM connectfour_games cg
            CROSS JOIN LATERAL jsonb_array_elements(
                CASE
                    WHEN jsonb_typeof(cg.game_state -> 'messages') = 'array' THEN cg.game_state -> 'messages'
                    ELSE '[]'::jsonb
                END
            ) AS msg(value)
            WHERE NOT EXISTS (
                SELECT 1
                FROM messages m
                WHERE m.game_id = cg.id
                  AND m.game_type = 'connectfour'
                  AND m.sender_id = msg.value ->> 'senderID'
                  AND m.sender_name = msg.value ->> 'senderName'
                  AND m.content = msg.value ->> 'content'
                  AND m.sent_at = COALESCE(NULLIF(msg.value ->> 'sentAt', '')::timestamptz, cg.finished_at)
            )
        $sql$;
    END IF;
END $$;

COMMIT;
