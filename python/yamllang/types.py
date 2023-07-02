from __future__ import annotations

import pydantic


class Module(pydantic.BaseModel):
    _path: str
    _package: Package

    @property
    def name(self) -> str:
        return self.path.split('.')[-1]

    @name.setter
    def name(self, name: str):
        self.path = self.path.rsplit('.', 1)[0] + '.' + name

    @property
    def path(self) -> str:
        return self._path

    @path.setter
    def path(self, path):
        path = '.'.join(path)
        if path != self._path:
            del self._package[self._path]
            self._path = '.'.join(path)
            self._package[path] = self

    @property
    def package(self):
        return self._package

    @package.setter
    def package(self, package: Package):
        if self._package is not package:
            self._package = package
            package[self.path] = self


class Package(pydantic.BaseModel):
    _path: str
    _modules_by_path: dict[str, Module]

    def __getitem__(self, path: str):
        return self._modules_by_path[path]

    def __setitem__(self, path: str, module: Module):
        if self._modules_by_path[path] is not module:
            module.package = self
            self._modules_by_path[path] = module

    def __delitem__(self, path: str):
        del self._modules_by_path[path]
