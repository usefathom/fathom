Ana Roadmap
===========

This is a general draft document for thoughts and todo's, without any structure to it.

### What's cooking?

- Return 0px GIF in /collect endpoint
- Hand out unique ID to each visitor
- Reference site URL when tracking.
- Reference path & title when tracking (indexed by path, update title when changes)
- Track referrals, use tables from aforementioned points.
- Bulk process tracking requests (Redis or in-memory?)
- Allow sorting in table overviews.
- Choose a OS license & settle on name.
- JS client for consuming API endpoints.
- Envelope API responses & perhaps return total in table overview?
- Track canonical URL's.
- Show referrals.
- Geolocate unknown IP addresses periodically.
- Mask last part of IP address.

### Key metrics

- Unique visits per day (in period)
- Pageviews per day (in period)
- Demographic
  - Country
  - Browser + version
  - Screen resolutions
- Acquisition
  - Referral's
  - Search keywords
