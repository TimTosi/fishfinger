CREATE DATABASE IF NOT EXISTS fishfinger;

CREATE TABLE IF NOT EXISTS fishfinger.sinking_ships (
    `id` int(10) unsigned NOT NULL,
    `name` varchar(50) NOT NULL
) ENGINE=MyISAM  DEFAULT CHARSET=utf8 ;

-- --------------------------------------------------------

--
-- Insert test data.
--

INSERT INTO fishfinger.sinking_ships
(`id`, `name`)
VALUES
(1, 'The Chimera'),
(2, 'The CMC'),
(3, 'The Last Immortal'),
(4, 'The Rageous Troll');
