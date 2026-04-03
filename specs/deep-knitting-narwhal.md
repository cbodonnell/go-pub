# go-pub UI Redesign Plan

## Context

The ActivityPub server (`go-pub`) is working well — federation, proxying, and content negotiation are solid. The existing frontend (`go-pub-client`) is a React 17 / Create React App / Bootstrap 4 app with class components and a garish hot-pink/cyan color palette. The goal is a complete frontend rewrite: modern React 18 + Vite, clean social-media aesthetic, and functional parity with the key features users love (feed, profiles, post composer, fediverse discoverability via click-through social graph).

The new client will live as a `client/` subdirectory inside the `go-pub` repo, so the whole project is self-contained. The Go server's `CLIENT` env var will point to the Vite build output at `./client/dist`.

---

## Tech Stack

| Concern | Choice | Why |
|---|---|---|
| Framework | React 18 + Vite | User preference; fast HMR, modern baseline |
| Language | TypeScript | Type-safe API shapes, better DX |
| Routing | React Router v6 | Modern file-based-friendly API |
| Data fetching | TanStack Query v5 | Infinite queries, caching, background refetch |
| Auth state | Zustand | Minimal global store for auth user |
| Styling | Tailwind CSS v3 | Utility-first, zero runtime, easy dark mode |
| Components | shadcn/ui | Unstyled accessible primitives, drop-in with Tailwind |
| HTTP | Axios (same pattern as existing) | Preserve proxy/cache interceptor logic |
| Icons | lucide-react | Clean icon set used by shadcn |

---

## Design System

**Color palette**: Neutral gray-based (slate family), one brand accent (indigo or sky blue). No hot pink.

**Typography**: Inter (via Fontsource). Clean sans-serif. Size scale: xs/sm/base/lg/xl.

**Mode**: Dark mode first, light mode toggle via `class` on `<html>` (Tailwind `darkMode: 'class'`).

**Card style**: Subtle border (`border border-slate-200 dark:border-slate-700`), rounded-xl, no harsh drop shadows.

**Layout**: Mastodon/Bluesky-inspired three-column desktop layout:
- Left: nav sidebar (sticky)
- Center: main content (scrolling feed/profiles)
- Right: optional sidebar (trending, discover)

Mobile: single column, bottom tab bar (5 tabs: Home, Explore, Compose, Notifications, Profile).

---

## App Architecture

```
go-pub/
  client/
    index.html
    vite.config.ts
    tailwind.config.ts
    src/
      main.tsx
      App.tsx              # Router setup
      api/
        client.ts          # Axios instances (proxy + cache interceptors)
        auth.ts            # Auth API calls
        actors.ts          # Actor/profile API calls
        activities.ts      # Activity API calls
        objects.ts         # Object/note API calls
        collections.ts     # OrderedCollection API calls
      types/
        activitypub.ts     # TypeScript types for AP models
      store/
        auth.ts            # Zustand auth store
      hooks/
        useActor.ts
        useCollection.ts
        useActivity.ts
        useObject.ts
        useFeed.ts
        useInteract.ts     # Like/follow/announce actions
      components/
        layout/
          AppLayout.tsx    # Three-column shell
          Sidebar.tsx      # Left nav
          Header.tsx       # Mobile top bar
          BottomNav.tsx    # Mobile bottom tab bar
        ui/                # shadcn primitives (Button, Card, Tabs, Avatar, Dialog, Toast)
        feed/
          FeedTimeline.tsx # Infinite scrolling feed
          PostCard.tsx     # Single post in the timeline
          PostDetail.tsx   # Full post with thread
          ReplyThread.tsx  # Nested reply chain
        composer/
          Composer.tsx     # Write/post modal (note + media)
        profile/
          ProfilePage.tsx  # Full actor profile
          ProfileHeader.tsx
          ProfileTabs.tsx
        discovery/
          ActorCard.tsx    # Actor in a list/grid
          CollectionBrowser.tsx  # Generic paginated collection
        auth/
          AuthModal.tsx    # Login / register dialog
        common/
          IriLink.tsx      # Smart link: local route or remote proxy route
          Skeleton.tsx     # Loading skeleton screens
          ErrorState.tsx
          EmptyState.tsx
      pages/
        HomePage.tsx       # / — feed for auth, discover for anon
        ProfilePage.tsx    # /users/:name
        PostPage.tsx       # /objects/:id
        ActivityPage.tsx   # /activities/:id
        RemoteActorPage.tsx    # /actor?remote=
        RemoteObjectPage.tsx   # /object?remote=
        RemoteCollectionPage.tsx  # /collection?remote=&type=
        NotFoundPage.tsx
```

---

## Key UX Decisions

### Feed / Timeline
- Infinite scroll via TanStack Query `useInfiniteQuery`
- Tabs above feed: **For You** (feed endpoint) | **Inbox** | **Outbox**
- Posts render as cards: avatar, name, handle, timestamp, content, attachment thumbnail, action bar (Like · Reply · Boost · ···)
- Action counts shown on buttons
- Clicking a post opens `PostPage` (thread view)

### Profiles
- Cover area (gradient fallback if no image) + large avatar
- Name, handle, follower/following counts as stats
- Follow / Message buttons for other users
- Tabs: **Posts** | **Replies** | **Liked**
- Followers/Following open as a sheet/drawer (not separate page), preserving discoverability inline

### Post Composer
- Floating button (bottom-right on desktop, center of bottom tab bar on mobile)
- Opens a modal: textarea with character count, media attach button
- Submit creates `Create > Note` activity via `POST /users/{name}/outbox`

### Social Graph Discoverability (preserve this!)
- Every actor name/avatar is a clickable `IriLink`
- Local IRIs → React Router navigation (no reload)
- Remote IRIs → `/actor?remote={encoded}` which proxies through the Go server's existing proxy infrastructure
- Same pattern for objects and collections
- "Browsing fediverse" feels seamless — same UI chrome for local and remote content

### Auth Gate
- Anon users see public profiles and posts
- Inbox/feed/interact require login — clicking those triggers `AuthModal`
- Auth via cookies + `withCredentials` (same as existing, no token storage change)

### Loading States
- Skeleton screens (not spinners) for profiles, post cards, collections
- Suspense boundaries at route level

### Notifications / Errors
- Sonner toast library for action feedback ("Liked!", "Followed!", error messages)
- No silent failures

---

## Implementation Phases

### Phase 1 — Project scaffold
- `npm create vite@latest client -- --template react-ts` inside `go-pub/`
- Install Tailwind, shadcn/ui init, React Router v6, TanStack Query, Zustand, Axios, lucide-react, sonner
- Configure Vite proxy for dev: forward `/users`, `/activities`, `/objects`, `/auth`, etc. to `http://localhost:8080`
- Port TypeScript types from existing `models.go` and AP spec
- Port `http.ts` / interceptors (proxy + cache logic) from go-pub-client

### Phase 2 — Auth + shell
- Zustand auth store + `useAuth` hook
- `AuthModal` component (login/register)
- `AppLayout` (sidebar + main + mobile bottom nav)
- `IriLink` component (smart routing)

### Phase 3 — Feed
- `useFeed` infinite query hook
- `PostCard` component
- `FeedTimeline` with infinite scroll
- `HomePage` with tab switcher (Feed / Inbox / Outbox)

### Phase 4 — Profiles
- `useActor` query hook
- `ProfileHeader`, `ProfileTabs`
- `ProfilePage` with Posts/Replies/Liked tabs
- `ActorCard` for use in follower/following lists
- Followers/Following drawer

### Phase 5 — Post detail + composer
- `PostPage` with thread view and replies
- `ReplyThread` recursive component
- `Composer` modal (create Note)
- `useInteract` hook (like, boost, follow, reply)
- Action bar on `PostCard`

### Phase 6 — Remote content / discovery
- `RemoteActorPage`, `RemoteObjectPage`, `RemoteCollectionPage`
- `CollectionBrowser` generic paginated browser
- Ensure all `IriLink` routes work end-to-end

### Phase 7 — Polish + deploy wiring
- Dark/light mode toggle
- Skeleton loaders everywhere
- Toast notifications
- Error/empty states
- `npm run build` output → set `CLIENT=./client/dist` in Go server `.env`
- Update `docker-compose` to run `npm run build` before serving

---

## Critical Files

**Backend (read-only reference):**
- `pkg/handlers/mux-handler.go` — all API routes
- `pkg/middleware/activitypub-middleware.go` — CLIENT proxy/static serving
- `pkg/models/models.go` — data shapes to type in TypeScript

**Existing client (reference, do not edit):**
- `~/code/go-pub-client/src/utils/http.js` — proxy + cache interceptors to port
- `~/code/go-pub-client/src/utils/urls.js` — IRI routing helpers to port
- `~/code/go-pub-client/src/utils/ap.js` — ActivityPub validation helpers to port
- `~/code/go-pub-client/src/components/` — all existing component logic to reference

**New files to create:**
- `go-pub/client/` — entire new frontend (Vite project)

---

## Verification

1. `cd client && npm run dev` — app runs at localhost:5173 with Vite proxy to Go server at :8080
2. Browse feed as authenticated user — posts load with infinite scroll
3. Click a local actor → profile page renders with tabs
4. Click a remote actor IRI → remote profile proxied and rendered in same UI
5. Compose a note → appears in outbox and feed
6. Like/follow → buttons update optimistically, API confirms
7. `npm run build` → `dist/` output; set `CLIENT=./client/dist` in Go server → serve via `go run ./cmd/pub serve`, verify HTML served on `/` and ActivityPub JSON on `/users/{name}` with correct `Accept` header
