-- DropForeignKey
ALTER TABLE "comentario" DROP CONSTRAINT "comentario_usuario_id_fkey";

-- DropForeignKey
ALTER TABLE "postagem" DROP CONSTRAINT "postagem_usuario_id_fkey";

-- AddForeignKey
ALTER TABLE "postagem" ADD CONSTRAINT "postagem_usuario_id_fkey" FOREIGN KEY ("usuario_id") REFERENCES "usuario"("id") ON DELETE SET NULL ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "comentario" ADD CONSTRAINT "comentario_usuario_id_fkey" FOREIGN KEY ("usuario_id") REFERENCES "usuario"("id") ON DELETE RESTRICT ON UPDATE CASCADE;
