package xml2map

import (
	"fmt"
	"github.com/google/gofuzz"
	"strings"
	"testing"
)

func BenchmarkDecoder(b *testing.B) {
	for n := 0; n < b.N; n++ {
		NewDecoder(strings.NewReader(`<container uid="FA6666D9-EC9F-4DA3-9C3D-4B2460A4E1F6" lifetime="2019-10-10T18:00:11">
				<cats>
					<cat>
						<id>CDA035B6-D453-4A17-B090-84295AE2DEC5</id>
						<name>moritz</name>
						<age>7</age> 	
						<items>
							<n>1293</n>
							<n>1255</n>
							<n>1257</n>
						</items>
					</cat>
					<cat>
						<id>1634C644-975F-4302-8336-1EF1366EC6A4</id>
						<name>oliver</name>
						<age>12</age>
					</cat>
				</cats>
				<color>white</color>
				<city>NY</city>
			</container>`)).Decode()

	}
}

func TestStartAttrs(t *testing.T) {
	tests := []string{
		`<container ="FA6666D9-EC9F-4DA3-9C3D-4B2460A4E1F6" lifetime="2019-10-10T18:00:11">
			<color>white</color>
		</container>`,
		`<container i=d="FA6666D9-EC9F-4DA3-9C3D-4B2460A4E1F6" lifetime="2019-10-10T18:00:11">
			<color>white</color>
		</container>`,
		`<container id="FA6666D9-EC9F-4DA3-9C3D-4B2460A4E1F6" lifetime="2019-10-10T18:00:11">
			<color id=>white</color>
		</container>`,
	}

	for _, s := range tests {
		_, err := NewDecoder(strings.NewReader(s)).Decode()
		if err == nil {
			t.Fail()
		}
	}
}

func TestNs(t *testing.T) {
	m, err := NewDecoder(strings.NewReader(
		`<container xmlns:h="http://www.w3.org/TR/html4/"
 						xmlns:xsl="http://www.w3.org/1999/XSL/Transform"></container>`)).Decode()
	if err != nil {
		t.Errorf("m: %v, err: %v\n", m, err)
	}
}

func TestPars(t *testing.T) {
	m, err := NewDecoder(strings.NewReader(
		`<customer  id="FA6666D9-EC9F-4DA3-9C3D-4B2460A4E1F6" lifetime="2019-10-10T18:00:11">
				<items id="100" count="3">
					<n id="10">1</n>
					<n id="20">2</n>
					<n id="30">3</n>
				</items>
			</customer>`)).Decode()

	if err != nil {
		t.Errorf("m: %v, err: %v\n", m, err)
	}
}

func TestFuzz1000(t *testing.T) {
	f := fuzz.New().NilChance(0).NumElements(1, 1000)
	var myMap map[string]int
	f.Fuzz(&myMap)

	for v := range myMap {
		m, err := NewDecoder(strings.NewReader(v)).Decode()
		if err == nil {
			fmt.Printf("m: %v", m)

		}
	}
}

func TestErrDecoder(t *testing.T) {
	m, err := NewDecoder(strings.NewReader(
		`<customer id="C1234">
			  <lname>Smith</lname>
			  <fname>John&gt;</fname>
			  <address type="biz">
				<street>1310 Villa Street</street>
				<city>Mountain View</city>
				<state>CA</state>
				<zip>94041</zip>
			  </address>
			<customer>`)).Decode()

	if m == nil && err != nil {
		t.Logf("result: %v err: %v\n", m, err)
	} else {
		t.Errorf("err %v\n", err)
	}
}

func TestEmpty(t *testing.T) {
	tests := []string{"", " ", "  ", ``, ` `, "\n"}

	for _, s := range tests {
		_, err := NewDecoder(strings.NewReader(s)).Decode()
		if err != ErrInvalidDocument {
			t.Fail()
		}
	}
}

func TestSpaces(t *testing.T) {
	m, err := NewDecoder(strings.NewReader(`   <note>
				  data
				</note>`)).Decode()

	if err != nil {
		t.Errorf("err %v\n", err)
	} else {
		if m["note"] != "data" {
			t.Errorf("data not found")
		}
	}
}

func TestInvalidStartIndex(t *testing.T) {
	_, err := NewDecoder(strings.NewReader(`d<note>
				  data
				</note>`)).Decode()

	if err == nil || err.Error() != "data at the root level is invalid" {
		t.Fail()
	}
}

func TestDecode(t *testing.T) {
	m, err := NewDecoder(strings.NewReader(
		`<container uid="FA6666D9-EC9F-4DA3-9C3D-4B2460A4E1F6" lifetime="2019-10-10T18:00:11">
				<cats>
					<cat>
						<id>CDA035B6-D453-4A17-B090-84295AE2DEC5</id>
						<name>moritz</name>
						<age>7</age> 	
						<items>
							<n>1293</n>
							<n>1255</n>
							<n>1257</n>
						</items>
					</cat>
					<cat>
						<id>1634C644-975F-4302-8336-1EF1366EC6A4</id>
						<name>oliver</name>
						<age>12</age>
						<items>
							<n>1293</n>
							<n>1255</n>
							<n>1257</n>
						</items>
					</cat>
					<dog color="gray">hello</dog>
				</cats>
				<color>white</color>
				<city>NY</city>
			</container>`)).Decode()

	if err != nil {
		t.Errorf("err: %v", err)
	}

	container := m["container"].(map[string]interface{})
	if container["@uid"] != "FA6666D9-EC9F-4DA3-9C3D-4B2460A4E1F6" && container["lifetime"] != "2019-10-10T18:00:11" {
		t.Errorf("container attrs not exists")
	} else {
		cats := container["cats"].(map[string]interface{})
		catsItems := cats["cat"].([]map[string]interface{})
		if len(catsItems) != 2 {
			t.Errorf("cats slice != 2")
		}

		dog := cats["dog"].(map[string]interface{})

		if dog["@color"] != "gray" || dog["#text"] != "hello" {
			t.Errorf("bad value or attr dog")
		}

		if container["color"] != "white" || container["city"] != "NY" {
			t.Errorf("bad value color")
		}

		cat := catsItems[0]
		if cat["id"] != "" && cat["name"] != "" && cat["age"] != "" {
			items := cat["items"].(map[string]interface{})["n"].([]string)
			if len(items) != 3 {
				t.Errorf("items len %v", len(items))
			}
		}
	}
}

func TestWithPrefix(t *testing.T) {
	m, err := NewDecoderWithPrefix(strings.NewReader(
		`<customer  id="FA6666D9-EC9F-4DA3-9C3D-4B2460A4E1F6" lifetime="2019-10-10T18:00:11">
				<items id="100" count="3">
					<n id="10">1</n>
					<n id="20">2</n>
					<n id="30">3</n>
				</items>
			</customer>`), "$", "#").Decode()

	if err != nil {
		t.Errorf("m: %v, err: %v\n", m, err)
	}

	customer := m["customer"].(map[string]interface{})
	if customer["$id"] != "FA6666D9-EC9F-4DA3-9C3D-4B2460A4E1F6" && customer["$lifetime"] != "2019-10-10T18:00:11" {
		t.Errorf("customer tag attr not found")
	} else {
		items := customer["items"].(map[string]interface{})
		if items["$id"] != "100" || items["$count"] != "3" {
			t.Errorf("items tag attr not found")
		} else {
			list := items["n"].([]map[string]interface{})
			if len(list) != 3 {
				t.Errorf("list len %v", len(items))
			} else {
				if list[1]["$id"] != "20" && list[1]["%"] != "2" {
					t.Errorf("invalid parse n element attr or text")
				}
			}
		}
	}
}

func TestWithNameSpace(t *testing.T) {
	m, err := NewDecoder(strings.NewReader(
		`<?xml version="1.0" encoding="utf-8"?>
		<rss version="2.0" xmlns:atom="http://www.w3.org/2005/Atom">
			<channel>
				<title>example.com RSS</title>
				<link>https://www.example.com/</link>
				<description>A cool website</description>
				<atom:link href="http://www.example.com/rss.xml" rel="self" type="application/rss+xml" />
				<atom:title>Atom Title</atom:title>
				<item>
					<title>Cool Article</title>
					<link>https://www.example.com/cool-article</link>
					<guid>https://www.example.com/cool-article</guid>
					<pubDate>Sun, 10 Dec 2017 05:00:00 GMT</pubDate>
					<description>My cool article description</description>
				</item>
			</channel>
		</rss>`)).Decode()

	if err != nil {
		t.Errorf("m: %v, err: %v\n", m, err)
	}

	rss := m["rss"].(map[string]interface{})["channel"].(map[string]interface{})
	if rss["atom:title"] != "Atom Title" {
		t.Errorf("invalid value for namespace node")
	}
}
