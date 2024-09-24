-- CreateTable
CREATE TABLE "usuarios" (
    "id" SERIAL NOT NULL,
    "login" VARCHAR(200) NOT NULL,
    "nome" VARCHAR(200) NOT NULL,
    "senha" VARCHAR(300) NOT NULL,

    CONSTRAINT "usuarios_pkey" PRIMARY KEY ("id")
);

-- CreateIndex
CREATE UNIQUE INDEX "usuarios_login_key" ON "usuarios"("login");
