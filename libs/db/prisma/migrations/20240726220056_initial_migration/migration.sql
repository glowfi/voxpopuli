-- CreateTable
CREATE TABLE "trophies" (
    "id" SERIAL NOT NULL,
    "title" TEXT NOT NULL,
    "description" TEXT NOT NULL,
    "imageLink" TEXT NOT NULL,

    CONSTRAINT "trophies_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "User" (
    "id" SERIAL NOT NULL,
    "username" TEXT NOT NULL,
    "email" TEXT NOT NULL,
    "password" TEXT NOT NULL,
    "cakeDay" INTEGER NOT NULL,
    "cakeDayHuman" TEXT NOT NULL,
    "accountAge" TEXT NOT NULL,
    "avatarImg" TEXT NOT NULL,
    "bannerImg" TEXT,
    "publicDescription" TEXT,
    "over18" BOOLEAN NOT NULL DEFAULT false,
    "suspended" BOOLEAN NOT NULL DEFAULT false,
    "keycolorHex" TEXT,
    "primarycolorHex" TEXT,
    "iconcolorHex" TEXT,
    "createdAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updatedAt" TIMESTAMP(3) NOT NULL,

    CONSTRAINT "User_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "awards" (
    "id" SERIAL NOT NULL,
    "title" TEXT NOT NULL,
    "imageLink" TEXT NOT NULL,

    CONSTRAINT "awards_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "topics" (
    "id" SERIAL NOT NULL,
    "title" TEXT NOT NULL,

    CONSTRAINT "topics_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "comments" (
    "id" SERIAL NOT NULL,
    "comment" TEXT NOT NULL,
    "commentHTML" TEXT NOT NULL,
    "authorId" INTEGER NOT NULL,
    "parentCommentId" INTEGER,
    "createdAtUnix" INTEGER NOT NULL,
    "createdAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updatedAt" TIMESTAMP(3) NOT NULL,
    "ups" INTEGER NOT NULL DEFAULT 0,
    "postId" INTEGER NOT NULL,

    CONSTRAINT "comments_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "mediaMetadata" (
    "id" SERIAL NOT NULL,
    "original_link" TEXT,
    "type" TEXT NOT NULL,
    "videoLink" TEXT,
    "imageLink" TEXT,
    "gifLink" TEXT,
    "gallery" JSONB,
    "postId" INTEGER NOT NULL,

    CONSTRAINT "mediaMetadata_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "flairs" (
    "id" SERIAL NOT NULL,
    "title" TEXT NOT NULL,
    "colorHex" TEXT,
    "voxsphereId" INTEGER NOT NULL,

    CONSTRAINT "flairs_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "rules" (
    "id" SERIAL NOT NULL,
    "title" TEXT NOT NULL,
    "description" TEXT,
    "voxsphereId" INTEGER NOT NULL,

    CONSTRAINT "rules_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "voxsphere" (
    "id" SERIAL NOT NULL,
    "name" TEXT NOT NULL,
    "about" TEXT,
    "logoURL" TEXT,
    "bannerURL" TEXT,
    "anchors" JSONB,
    "buttonColorHex" TEXT,
    "headerColorHex" TEXT,
    "bannerBgColorHex" TEXT,
    "createdAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updatedAt" TIMESTAMP(3) NOT NULL,
    "createdAtUnix" INTEGER NOT NULL,
    "createdatHuman" TEXT NOT NULL,
    "totalmembers" INTEGER NOT NULL,
    "totalmembersHuman" TEXT NOT NULL,
    "over18" BOOLEAN NOT NULL DEFAULT false,
    "spoilersEnabled" BOOLEAN NOT NULL DEFAULT false,

    CONSTRAINT "voxsphere_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "posts" (
    "id" SERIAL NOT NULL,
    "authorId" INTEGER NOT NULL,
    "createdAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updatedAt" TIMESTAMP(3) NOT NULL,
    "createdAtUnix" INTEGER NOT NULL,
    "createdatHuman" TEXT NOT NULL,
    "linkType" BOOLEAN NOT NULL DEFAULT false,
    "numComments" INTEGER NOT NULL DEFAULT 0,
    "over18" BOOLEAN NOT NULL,
    "spoiler" BOOLEAN NOT NULL,
    "title" TEXT NOT NULL,
    "text" TEXT NOT NULL,
    "textHTML" TEXT,
    "ups" INTEGER NOT NULL,

    CONSTRAINT "posts_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "_voxsphere_member" (
    "A" INTEGER NOT NULL,
    "B" INTEGER NOT NULL
);

-- CreateTable
CREATE TABLE "_voxsphere_moderator" (
    "A" INTEGER NOT NULL,
    "B" INTEGER NOT NULL
);

-- CreateTable
CREATE TABLE "_UserTotrophies" (
    "A" INTEGER NOT NULL,
    "B" INTEGER NOT NULL
);

-- CreateTable
CREATE TABLE "_awardsToposts" (
    "A" INTEGER NOT NULL,
    "B" INTEGER NOT NULL
);

-- CreateTable
CREATE TABLE "_topicsTovoxsphere" (
    "A" INTEGER NOT NULL,
    "B" INTEGER NOT NULL
);

-- CreateTable
CREATE TABLE "_flairsToposts" (
    "A" INTEGER NOT NULL,
    "B" INTEGER NOT NULL
);

-- CreateTable
CREATE TABLE "_postsTotopics" (
    "A" INTEGER NOT NULL,
    "B" INTEGER NOT NULL
);

-- CreateTable
CREATE TABLE "_postsTovoxsphere" (
    "A" INTEGER NOT NULL,
    "B" INTEGER NOT NULL
);

-- CreateIndex
CREATE UNIQUE INDEX "trophies_title_key" ON "trophies"("title");

-- CreateIndex
CREATE UNIQUE INDEX "User_username_key" ON "User"("username");

-- CreateIndex
CREATE UNIQUE INDEX "User_email_key" ON "User"("email");

-- CreateIndex
CREATE UNIQUE INDEX "awards_title_key" ON "awards"("title");

-- CreateIndex
CREATE UNIQUE INDEX "topics_title_key" ON "topics"("title");

-- CreateIndex
CREATE UNIQUE INDEX "voxsphere_name_key" ON "voxsphere"("name");

-- CreateIndex
CREATE UNIQUE INDEX "_voxsphere_member_AB_unique" ON "_voxsphere_member"("A", "B");

-- CreateIndex
CREATE INDEX "_voxsphere_member_B_index" ON "_voxsphere_member"("B");

-- CreateIndex
CREATE UNIQUE INDEX "_voxsphere_moderator_AB_unique" ON "_voxsphere_moderator"("A", "B");

-- CreateIndex
CREATE INDEX "_voxsphere_moderator_B_index" ON "_voxsphere_moderator"("B");

-- CreateIndex
CREATE UNIQUE INDEX "_UserTotrophies_AB_unique" ON "_UserTotrophies"("A", "B");

-- CreateIndex
CREATE INDEX "_UserTotrophies_B_index" ON "_UserTotrophies"("B");

-- CreateIndex
CREATE UNIQUE INDEX "_awardsToposts_AB_unique" ON "_awardsToposts"("A", "B");

-- CreateIndex
CREATE INDEX "_awardsToposts_B_index" ON "_awardsToposts"("B");

-- CreateIndex
CREATE UNIQUE INDEX "_topicsTovoxsphere_AB_unique" ON "_topicsTovoxsphere"("A", "B");

-- CreateIndex
CREATE INDEX "_topicsTovoxsphere_B_index" ON "_topicsTovoxsphere"("B");

-- CreateIndex
CREATE UNIQUE INDEX "_flairsToposts_AB_unique" ON "_flairsToposts"("A", "B");

-- CreateIndex
CREATE INDEX "_flairsToposts_B_index" ON "_flairsToposts"("B");

-- CreateIndex
CREATE UNIQUE INDEX "_postsTotopics_AB_unique" ON "_postsTotopics"("A", "B");

-- CreateIndex
CREATE INDEX "_postsTotopics_B_index" ON "_postsTotopics"("B");

-- CreateIndex
CREATE UNIQUE INDEX "_postsTovoxsphere_AB_unique" ON "_postsTovoxsphere"("A", "B");

-- CreateIndex
CREATE INDEX "_postsTovoxsphere_B_index" ON "_postsTovoxsphere"("B");

-- AddForeignKey
ALTER TABLE "comments" ADD CONSTRAINT "comments_authorId_fkey" FOREIGN KEY ("authorId") REFERENCES "User"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "comments" ADD CONSTRAINT "comments_parentCommentId_fkey" FOREIGN KEY ("parentCommentId") REFERENCES "comments"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "comments" ADD CONSTRAINT "comments_postId_fkey" FOREIGN KEY ("postId") REFERENCES "posts"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "mediaMetadata" ADD CONSTRAINT "mediaMetadata_postId_fkey" FOREIGN KEY ("postId") REFERENCES "posts"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "flairs" ADD CONSTRAINT "flairs_voxsphereId_fkey" FOREIGN KEY ("voxsphereId") REFERENCES "voxsphere"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "rules" ADD CONSTRAINT "rules_voxsphereId_fkey" FOREIGN KEY ("voxsphereId") REFERENCES "voxsphere"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "posts" ADD CONSTRAINT "posts_authorId_fkey" FOREIGN KEY ("authorId") REFERENCES "User"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "_voxsphere_member" ADD CONSTRAINT "_voxsphere_member_A_fkey" FOREIGN KEY ("A") REFERENCES "User"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "_voxsphere_member" ADD CONSTRAINT "_voxsphere_member_B_fkey" FOREIGN KEY ("B") REFERENCES "voxsphere"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "_voxsphere_moderator" ADD CONSTRAINT "_voxsphere_moderator_A_fkey" FOREIGN KEY ("A") REFERENCES "User"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "_voxsphere_moderator" ADD CONSTRAINT "_voxsphere_moderator_B_fkey" FOREIGN KEY ("B") REFERENCES "voxsphere"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "_UserTotrophies" ADD CONSTRAINT "_UserTotrophies_A_fkey" FOREIGN KEY ("A") REFERENCES "User"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "_UserTotrophies" ADD CONSTRAINT "_UserTotrophies_B_fkey" FOREIGN KEY ("B") REFERENCES "trophies"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "_awardsToposts" ADD CONSTRAINT "_awardsToposts_A_fkey" FOREIGN KEY ("A") REFERENCES "awards"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "_awardsToposts" ADD CONSTRAINT "_awardsToposts_B_fkey" FOREIGN KEY ("B") REFERENCES "posts"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "_topicsTovoxsphere" ADD CONSTRAINT "_topicsTovoxsphere_A_fkey" FOREIGN KEY ("A") REFERENCES "topics"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "_topicsTovoxsphere" ADD CONSTRAINT "_topicsTovoxsphere_B_fkey" FOREIGN KEY ("B") REFERENCES "voxsphere"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "_flairsToposts" ADD CONSTRAINT "_flairsToposts_A_fkey" FOREIGN KEY ("A") REFERENCES "flairs"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "_flairsToposts" ADD CONSTRAINT "_flairsToposts_B_fkey" FOREIGN KEY ("B") REFERENCES "posts"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "_postsTotopics" ADD CONSTRAINT "_postsTotopics_A_fkey" FOREIGN KEY ("A") REFERENCES "posts"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "_postsTotopics" ADD CONSTRAINT "_postsTotopics_B_fkey" FOREIGN KEY ("B") REFERENCES "topics"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "_postsTovoxsphere" ADD CONSTRAINT "_postsTovoxsphere_A_fkey" FOREIGN KEY ("A") REFERENCES "posts"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "_postsTovoxsphere" ADD CONSTRAINT "_postsTovoxsphere_B_fkey" FOREIGN KEY ("B") REFERENCES "voxsphere"("id") ON DELETE CASCADE ON UPDATE CASCADE;
