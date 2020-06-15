package search

import (
	"testing"
	"time"

	"github.com/davecgh/go-spew/spew"

	"github.com/google/go-cmp/cmp"
	"github.com/hanzki/moviebox-server/core"
)

const sample = `
<rss version="1.0" xmlns:atom="http://www.w3.org/2005/Atom" xmlns:torznab="http://torznab.com/schemas/2015/feed">
	<channel>
		<atom:link href="http://localhost:9117/" rel="self" type="application/rss+xml" />
		<title>AggregateSearch</title>
		<description>This feed includes all configured trackers</description>
		<link>http://127.0.0.1/</link>
		<language>en-us</language>
		<category>search</category>
		<image>
			<url>http://localhost:9117/logos/all.png</url>
			<title>AggregateSearch</title>
			<link>http://127.0.0.1/</link>
			<description>AggregateSearch</description>
		</image>
        <item>
            <title>Big Buck Bunny (1280x720 msmp4)</title>
            <guid>http://www.legittorrents.info/index.php?page=torrent-details&amp;id=d42644b1a0eb635a4ffaf6ce17042d71428e6a8e</guid>
            <jackettindexer id="legittorrents">Legit Torrents</jackettindexer>
            <comments>http://www.legittorrents.info/index.php?page=torrent-details&amp;id=d42644b1a0eb635a4ffaf6ce17042d71428e6a8e</comments>
            <pubDate>Sun, 01 Jun 2008 00:00:00 +0000</pubDate>
            <grabs>7940</grabs>
			<description />
            <link>http://localhost:9117/dl/legittorrents/?jackett_apikey=gpowobdo7ztigmxjoeokamhjjh7bz8us&amp;path=Q2ZESjhPcWNEVnJwNmdWQmhqNDh0dGcxM2VtcXRJM0c3WGhuOUZCZ3lSeVRyWjkxWGhheThGd3VMWUc3RnM0d0xvWkhHMk5meDRRY3VCRk02TlExSjF5Tl81X0thTWY5ZzlWSE9xU1I5SlBIdjk2MEJIbXNHSElKdXJRMFdFVF81M1hETGNkaFM2MmFRaTU3aTczT0hsUm8zRE1zd2hRY1AtcHd5a1prSTNzdDhRcWp4bm5lTm1CVEREdEdSb0MyOEE0OEY1WnZZVTYwOW4zWWw4Zkh4QXpJa3BpbGtGT1VDa2FqcnNBWVlpRGY3N3IzLWdQQ1JKc0dkNld2TGlIbzNQLUZtcUdHcDYyRkdyVVNVZmplV1dWc3IyeDNhLXA1MjhVNXRmQXZZbk5OeERHdg&amp;file=Big+Buck+Bunny+(1280x720+msmp4)</link>
            <category>2000</category>
            <category>100001</category>
            <enclosure url="http://localhost:9117/dl/legittorrents/?jackett_apikey=gpowobdo7ztigmxjoeokamhjjh7bz8us&amp;path=Q2ZESjhPcWNEVnJwNmdWQmhqNDh0dGcxM2VtcXRJM0c3WGhuOUZCZ3lSeVRyWjkxWGhheThGd3VMWUc3RnM0d0xvWkhHMk5meDRRY3VCRk02TlExSjF5Tl81X0thTWY5ZzlWSE9xU1I5SlBIdjk2MEJIbXNHSElKdXJRMFdFVF81M1hETGNkaFM2MmFRaTU3aTczT0hsUm8zRE1zd2hRY1AtcHd5a1prSTNzdDhRcWp4bm5lTm1CVEREdEdSb0MyOEE0OEY1WnZZVTYwOW4zWWw4Zkh4QXpJa3BpbGtGT1VDa2FqcnNBWVlpRGY3N3IzLWdQQ1JKc0dkNld2TGlIbzNQLUZtcUdHcDYyRkdyVVNVZmplV1dWc3IyeDNhLXA1MjhVNXRmQXZZbk5OeERHdg&amp;file=Big+Buck+Bunny+(1280x720+msmp4)" type="application/x-bittorrent" />
            <torznab:attr name="category" value="2000" />
            <torznab:attr name="category" value="100001" />
            <torznab:attr name="seeders" value="59" />
            <torznab:attr name="peers" value="59" />
            <torznab:attr name="minimumratio" value="1" />
            <torznab:attr name="minimumseedtime" value="172800" />
            <torznab:attr name="downloadvolumefactor" value="0" />
            <torznab:attr name="uploadvolumefactor" value="1" />
        </item>
	</channel>
</rss>
`

func TestParseTorznab(t *testing.T) {
	pubDate, err := time.Parse(time.RFC1123Z, "Sun, 01 Jun 2008 00:00:00 +0000")
	if err != nil {
		panic(err)
	}

	want := &core.SearchResult{
		Title:       "Big Buck Bunny (1280x720 msmp4)",
		IndexerGUID: "http://www.legittorrents.info/index.php?page=torrent-details&id=d42644b1a0eb635a4ffaf6ce17042d71428e6a8e",
		Indexer:     "legittorrents",
		PubDate:     pubDate,
		Link:        "http://localhost:9117/dl/legittorrents/?jackett_apikey=gpowobdo7ztigmxjoeokamhjjh7bz8us&path=Q2ZESjhPcWNEVnJwNmdWQmhqNDh0dGcxM2VtcXRJM0c3WGhuOUZCZ3lSeVRyWjkxWGhheThGd3VMWUc3RnM0d0xvWkhHMk5meDRRY3VCRk02TlExSjF5Tl81X0thTWY5ZzlWSE9xU1I5SlBIdjk2MEJIbXNHSElKdXJRMFdFVF81M1hETGNkaFM2MmFRaTU3aTczT0hsUm8zRE1zd2hRY1AtcHd5a1prSTNzdDhRcWp4bm5lTm1CVEREdEdSb0MyOEE0OEY1WnZZVTYwOW4zWWw4Zkh4QXpJa3BpbGtGT1VDa2FqcnNBWVlpRGY3N3IzLWdQQ1JKc0dkNld2TGlIbzNQLUZtcUdHcDYyRkdyVVNVZmplV1dWc3IyeDNhLXA1MjhVNXRmQXZZbk5OeERHdg&file=Big+Buck+Bunny+(1280x720+msmp4)",
		Categories:  []string{"Movies"},
		Seeds:       59,
		Peers:       59,
	}

	categories := map[string]*JackettCategory{"2000": {"2000", "Movies", nil}}
	results := parseTorznab([]byte(sample), categories)

	if len(results) != 1 {
		t.Errorf("parseTorznab: Got wrong number of results. Expected 1, got %d", len(results))
		return
	}

	if got := results[0]; !cmp.Equal(got, want) {
		t.Errorf("parseTorznab:\nExpected %+v\nGot %+v", want, got)
		spew.Dump(got, want)
	}
}

const capsSample = `
<?xml version="1.0" encoding="UTF-8"?>
<caps>
    <server title="Jackett" />
    <limits max="1000" default="1000" />
    <searching>
        <search available="yes" supportedParams="q" />
        <tv-search available="yes" supportedParams="q,season,ep" />
        <movie-search available="yes" supportedParams="q" />
        <music-search available="no" supportedParams="" />
        <audio-search available="no" supportedParams="" />
    </searching>
    <categories>
        <category id="1000" name="Console">
            <subcat id="1010" name="Console/NDS" />
            <subcat id="1020" name="Console/PSP" />
            <subcat id="1030" name="Console/Wii" />
            <subcat id="1040" name="Console/Xbox" />
            <subcat id="1050" name="Console/Xbox 360" />
            <subcat id="1060" name="Console/Wiiware/VC" />
            <subcat id="1070" name="Console/XBOX 360 DLC" />
            <subcat id="1080" name="Console/PS3" />
            <subcat id="1090" name="Console/Other" />
            <subcat id="1110" name="Console/3DS" />
            <subcat id="1120" name="Console/PS Vita" />
            <subcat id="1130" name="Console/WiiU" />
            <subcat id="1140" name="Console/Xbox One" />
            <subcat id="1180" name="Console/PS4" />
        </category>
        <category id="1010" name="Console/NDS" />
        <category id="1020" name="Console/PSP" />
        <category id="1030" name="Console/Wii" />
        <category id="1040" name="Console/Xbox" />
        <category id="1050" name="Console/Xbox 360" />
        <category id="1080" name="Console/PS3" />
        <category id="1090" name="Console/Other" />
        <category id="1110" name="Console/3DS" />
        <category id="1180" name="Console/PS4" />
        <category id="2000" name="Movies">
            <subcat id="2010" name="Movies/Foreign" />
            <subcat id="2020" name="Movies/Other" />
            <subcat id="2030" name="Movies/SD" />
            <subcat id="2040" name="Movies/HD" />
            <subcat id="2045" name="Movies/UHD" />
            <subcat id="2050" name="Movies/3D" />
            <subcat id="2060" name="Movies/BluRay" />
            <subcat id="2070" name="Movies/DVD" />
            <subcat id="2080" name="Movies/WEBDL" />
        </category>
        <category id="2010" name="Movies/Foreign" />
        <category id="2020" name="Movies/Other" />
        <category id="2030" name="Movies/SD" />
        <category id="2040" name="Movies/HD" />
        <category id="2045" name="Movies/UHD" />
        <category id="2050" name="Movies/3D" />
        <category id="2060" name="Movies/BluRay" />
        <category id="2070" name="Movies/DVD" />
        <category id="3000" name="Audio">
            <subcat id="3010" name="Audio/MP3" />
            <subcat id="3020" name="Audio/Video" />
            <subcat id="3030" name="Audio/Audiobook" />
            <subcat id="3040" name="Audio/Lossless" />
            <subcat id="3050" name="Audio/Other" />
            <subcat id="3060" name="Audio/Foreign" />
        </category>
        <category id="3010" name="Audio/MP3" />
        <category id="3020" name="Audio/Video" />
        <category id="3030" name="Audio/Audiobook" />
        <category id="3040" name="Audio/Lossless" />
        <category id="3050" name="Audio/Other" />
        <category id="4000" name="PC">
            <subcat id="4010" name="PC/0day" />
            <subcat id="4020" name="PC/ISO" />
            <subcat id="4030" name="PC/Mac" />
            <subcat id="4040" name="PC/Phone-Other" />
            <subcat id="4050" name="PC/Games" />
            <subcat id="4060" name="PC/Phone-IOS" />
            <subcat id="4070" name="PC/Phone-Android" />
        </category>
        <category id="4010" name="PC/0day" />
        <category id="4030" name="PC/Mac" />
        <category id="4040" name="PC/Phone-Other" />
        <category id="4050" name="PC/Games" />
        <category id="4060" name="PC/Phone-IOS" />
        <category id="4070" name="PC/Phone-Android" />
        <category id="5000" name="TV">
            <subcat id="5010" name="TV/WEB-DL" />
            <subcat id="5020" name="TV/FOREIGN" />
            <subcat id="5030" name="TV/SD" />
            <subcat id="5040" name="TV/HD" />
            <subcat id="5045" name="TV/UHD" />
            <subcat id="5050" name="TV/OTHER" />
            <subcat id="5060" name="TV/Sport" />
            <subcat id="5070" name="TV/Anime" />
            <subcat id="5080" name="TV/Documentary" />
        </category>
        <category id="5030" name="TV/SD" />
        <category id="5040" name="TV/HD" />
        <category id="5045" name="TV/UHD" />
        <category id="5050" name="TV/OTHER" />
        <category id="5060" name="TV/Sport" />
        <category id="5070" name="TV/Anime" />
        <category id="5080" name="TV/Documentary" />
        <category id="6000" name="XXX">
            <subcat id="6010" name="XXX/DVD" />
            <subcat id="6020" name="XXX/WMV" />
            <subcat id="6030" name="XXX/XviD" />
            <subcat id="6040" name="XXX/x264" />
            <subcat id="6050" name="XXX/Other" />
            <subcat id="6060" name="XXX/Imageset" />
            <subcat id="6070" name="XXX/Packs" />
        </category>
        <category id="6010" name="XXX/DVD" />
        <category id="6040" name="XXX/x264" />
        <category id="6060" name="XXX/Imageset" />
        <category id="7000" name="Other">
            <subcat id="7010" name="Other/Misc" />
            <subcat id="7020" name="Other/Hashed" />
        </category>
        <category id="7010" name="Other/Misc" />
        <category id="8000" name="Books">
            <subcat id="8010" name="Books/Ebook" />
            <subcat id="8020" name="Books/Comics" />
            <subcat id="8030" name="Books/Magazines" />
            <subcat id="8040" name="Books/Technical" />
            <subcat id="8050" name="Books/Other" />
            <subcat id="8060" name="Books/Foreign" />
        </category>
        <category id="8010" name="Books/Ebook" />
        <category id="8020" name="Books/Comics" />
        <category id="8030" name="Books/Magazines" />
        <category id="8040" name="Books/Technical" />
        <category id="8050" name="Books/Other" />
    </categories>
</caps>
`

func TestParseCaps(t *testing.T) {

	want := &JackettCategory{
		ID:   "5000",
		Name: "TV",
		SubCategories: []*JackettCategory{
			&JackettCategory{"5010", "TV/WEB-DL", nil},
			&JackettCategory{"5020", "TV/FOREIGN", nil},
			&JackettCategory{"5030", "TV/SD", nil},
			&JackettCategory{"5040", "TV/HD", nil},
			&JackettCategory{"5045", "TV/UHD", nil},
			&JackettCategory{"5050", "TV/OTHER", nil},
			&JackettCategory{"5060", "TV/Sport", nil},
			&JackettCategory{"5070", "TV/Anime", nil},
			&JackettCategory{"5080", "TV/Documentary", nil},
		},
	}
	results := parseCaps([]byte(capsSample))

	var got *JackettCategory
	for _, v := range results {
		if v.ID == want.ID {
			got = v
			break
		}
	}

	if !cmp.Equal(got, want) {
		t.Errorf("parseCaps:\nExpected %+v\nGot %+v", want, got)
		spew.Dump(got, want)
	}
}
